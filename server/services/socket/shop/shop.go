package shop

import (
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"main/server/services/socket"
	"main/server/utils"

	socketio "github.com/googollee/go-socket.io"
)

func BuyFromShop(s socketio.Conn, reqData map[string]interface{}) {

	playerId := s.Context().(string)
	buyId, ok := reqData["id"].(string)
	if !ok {
		response.SocketResponse("Id is required", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "buyStore", s)
		return
	}
	// fmt.Println(buyId)
	var temp model.Shop
	query := "SELECT * FROM shops WHERE id=?"
	err := db.QueryExecutor(query, &temp, buyId)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "buyStore", s)
		return
	}
	// fmt.Println("temp is ", temp)

	playerDetails, err := utils.GetPlayerDetails(playerId)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "buyStore", s)
		return
	}

	//Check if the player if eligible for buying the things
	switch temp.PurchaseType {
	case utils.COINS:
		{
			if temp.PurchaseValue > playerDetails.Coins {
				response.SocketResponse(utils.NOT_ENOUGH_COINS, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "buyStore", s)
				return
			}
		}
	case utils.CASH:
		{
			if temp.PurchaseValue > playerDetails.Cash {

				response.SocketResponse(utils.NOT_ENOUGH_CASH, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "buyStore", s)
				return
			}
		}
	case utils.REPAIR_PARTS:
		{
			if temp.PurchaseValue > playerDetails.RepairCurrency {
				response.SocketResponse(utils.NOT_ENOUGH_REPAIR_PARTS, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "buyStore", s)
				return
			}
		}
	}

	switch temp.PurchaseType {
	case utils.COINS:
		playerDetails.Coins -= temp.PurchaseValue
	case utils.CASH:
		playerDetails.Cash -= temp.PurchaseValue
	case utils.REPAIR_PARTS:
		playerDetails.RepairCurrency -= temp.PurchaseValue
	}

	switch temp.RewardType {
	case utils.COINS:
		playerDetails.Coins += temp.RewardValue
	case utils.CASH:
		playerDetails.Cash += temp.RewardValue
	case utils.REPAIR_PARTS:
		playerDetails.RepairCurrency += temp.RewardValue
	}

	err = db.UpdateRecord(&playerDetails, playerId, "player_id").Error
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "buyStore", s)
		return
	}

	playerResponse, err := socket.GetPlayerDetailsCopy(playerDetails.PlayerId)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "buyStore", s)
		return
	}

	utils.SocketServerInstance.BroadcastToRoom("/", playerId, "playerDetails", playerResponse)

	response.SocketResponse(utils.PURCHASE_SUCCESS, utils.HTTP_OK, utils.SUCCESS, temp, "buyStore", s)

}
