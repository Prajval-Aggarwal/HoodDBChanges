package car

import (
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"main/server/utils"

	socketio "github.com/googollee/go-socket.io"
)

func BuyCarService(s socketio.Conn, req map[string]interface{}) {
	fmt.Println("Buy car socket handler called")
	playerId := s.Context().(string)
	carId, ok := req["carId"].(string)
	if !ok {
		response.SocketResponse("Car id is required", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "carBuy", s)
		return
	}
	purchaseType, ok := req["purchaseType"].(float64)
	if !ok {
		response.SocketResponse("Purchase type is required", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "carBuy", s)
		return
	}

	// Check if the car exists in the database.
	if !db.RecordExist("cars", carId, "car_id") {
		response.SocketResponse(utils.CAR_NOT_FOUND, utils.HTTP_NOT_FOUND, utils.FAILURE, nil, "carBuy", s)
		return
	}

	// Fetch car details from the database.
	var carDetails model.Car
	err := db.FindById(&carDetails, carId, "car_id")
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "carBuy", s)
		return
	}

	// Fetch player details from the database.

	playerDetails, err := utils.GetPlayerDetails(playerId)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "carBuy", s)
		return
	}
	//mapping purchaseType
	var buyType string
	switch purchaseType {
	case utils.COINS:
		buyType = "coins"
	case utils.CASH:
		buyType = "cash"
	}

	var amount int64

	// Check if the player has enough currency to buy the car.

	if carDetails.CurrType == "coins" {
		amount = playerDetails.Coins
	} else {
		amount = playerDetails.Cash
	}

	if amount < int64(carDetails.CurrAmount) {
		response.SocketResponse(utils.NOT_ENOUGH_COINS, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "carBuy", s)
		return
	}

	// Start a database transaction to handle the purchase.
	tx := db.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if buyType == carDetails.CurrType {
		if carDetails.CurrType == "coins" {
			playerDetails.Coins -= int64(carDetails.CurrAmount)
		} else {
			playerDetails.Cash -= int64(carDetails.CurrAmount)
		}
	} else {
		playerDetails.Cash -= int64(carDetails.PremiumBuy)

	}

	//updating player details
	err = tx.Where(utils.PLAYER_ID, playerId).Updates(&playerDetails).Error
	if err != nil {
		tx.Rollback()
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "carBuy", s)

		return
	}

	//Unequip the currenntluy selected car and euip the recently bought car
	query := "UPDATE owned_cars SET selected = false WHERE player_id=? AND selected=true"
	err = tx.Exec(query, playerId).Error
	if err != nil {
		tx.Rollback()
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "carBuy", s)

		return
	}

	//finding the default customisation for that car
	//Adding these default to player_car_customisations
	//Also give a garage to that person
	err = utils.SetCarData(carId, playerId)
	if err != nil {
		tx.Rollback()
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "carBuy", s)
		return
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		response.SocketResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "carBuy", s)
		return
	}
	response.SocketResponse(utils.CAR_BOUGHT_SUCESS, utils.HTTP_OK, utils.SUCCESS, carDetails, "carBuy", s)

}
