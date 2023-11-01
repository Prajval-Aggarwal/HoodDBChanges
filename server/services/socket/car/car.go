package car

import (
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"main/server/services/socket"
	"main/server/utils"
	"math"

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

	// Start a database transaction to handle the purchase.

	if buyType == carDetails.CurrType {
		if carDetails.CurrType == "coins" {
			if carDetails.CurrAmount > float64(playerDetails.Coins) {
				response.SocketResponse(utils.NOT_ENOUGH_COINS, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "carBuy", s)
				return
			} else {
				playerDetails.Coins -= int64(carDetails.CurrAmount)
			}
		} else {
			if carDetails.CurrAmount > float64(playerDetails.Cash) {
				response.SocketResponse(utils.NOT_ENOUGH_CASH, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "carBuy", s)
				return
			} else {
				playerDetails.Cash -= int64(carDetails.CurrAmount)
			}
		}
	} else {
		if float64(carDetails.PremiumBuy) > float64(playerDetails.Cash) {
			response.SocketResponse(utils.NOT_ENOUGH_CASH, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "carBuy", s)
			return
		} else {
			playerDetails.Cash -= int64(carDetails.PremiumBuy)
		}
	}
	fmt.Println("Player details is", playerDetails)

	//updating player details
	err = db.UpdateRecord(&playerDetails, playerId, "player_id").Error
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "carBuy", s)

		return
	}

	//Unequip the currenntluy selected car and euip the recently bought car
	query := "UPDATE owned_cars SET selected = false WHERE player_id=? AND selected=true"
	err = db.RawExecutor(query, playerId)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "carBuy", s)
		return
	}
	fmt.Println("asvdjadjsajd")

	//finding the default customisation for that car
	//Adding these default to player_car_customisations
	//Also give a garage to that person
	err = utils.SetCarData(carId, playerId)
	if err != nil {
		fmt.Println("error is", err)
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "carBuy", s)
		return
	}
	fmt.Println("ajsvdjajshlalvalv")

	playerResponse, err := socket.GetPlayerDetailsCopy(playerDetails.PlayerId)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "carBuy", s)
		return
	}

	fmt.Println("Player Response", *playerResponse)
	//braodcasting the updated player details to the front end

	if !utils.SocketServerInstance.BroadcastToRoom("/", playerId, "playerDetails", *playerResponse) {
		fmt.Println("advjabdjkjasd")
		return
	}
	response.SocketResponse(utils.CAR_BOUGHT_SUCESS, utils.HTTP_OK, utils.SUCCESS, carDetails, "carBuy", s)

}

func RepairCarService(s socketio.Conn, reqData map[string]interface{}) {

	fmt.Println("Repair car socket handler called...")
	playerId := s.Context().(string)
	custId, ok := reqData["custId"].(string)
	if !ok {
		response.SocketResponse(utils.CUSTID_REQUIRED, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "carRepair", s)
		return
	}

	playerDetails, err := utils.GetPlayerDetails(playerId)
	if err != nil {

		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "carRepair", s)
		return
	}

	// Fetch car stats details from the database.
	var carDetails model.DefaultCustomisation

	query := "SELECT * FROM default_customisation dc JOIN cars c on c.car_id=dc.car_id WHERE dc.cust_id=?"
	err = db.QueryExecutor(query, &carDetails, custId)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "carRepair", s)
		return
	}

	// Fetch player car stats details from the database.
	var playerCarStats model.PlayerCarCustomisation
	query = "SELECT * FROM player_cars_stats WHERE cust_id=? AND player_id=?"
	err = db.QueryExecutor(query, &playerCarStats, custId, playerId)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "carRepair", s)
		return
	}

	// Calculate the difference in durability to repair.
	durabilityDiff := carDetails.Durability - playerCarStats.Durability

	// Check if the player has enough repair parts to perform the repair.
	if int64(math.Floor(float64(durabilityDiff)*2.5)) > playerDetails.RepairCurrency {
		response.SocketResponse(utils.NOT_ENOUGH_REPAIR_PARTS, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "carRepair", s)
		return
	}

	// Start a database transaction to handle the car repair.
	tx := db.BeginTransaction()
	if tx.Error != nil {
		response.SocketResponse(tx.Error.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "carRepair", s)
		return
	}
	// Update the player car's durability to the car's maximum durability.
	playerCarStats.Durability = carDetails.Durability
	query = "UPDATE player_car_customisations SET durability = ? WHERE player_id=? AND cust_id=?"
	err = tx.Exec(query, carDetails.Durability, playerId, custId).Error
	if err != nil {
		tx.Rollback() // Rollback the transaction if there's an error.
		response.SocketResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "carRepair", s)
		return
	}
	if err := tx.Commit().Error; err != nil {
		response.SocketResponse(tx.Error.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "carRepair", s)
		return
	}

	utils.SocketServerInstance.BroadcastToRoom("/", playerId, "carDetails", carDetails)
	response.SocketResponse(utils.CAR_REPAIR_SUCCESS, utils.HTTP_OK, utils.SUCCESS, nil, "carRepair", s)
}
