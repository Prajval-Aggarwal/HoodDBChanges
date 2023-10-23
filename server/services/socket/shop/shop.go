package shop

import (
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/response"
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
	fmt.Println(buyId)
	var temp model.Shop
	query := "SELECT * FROM shops WHERE id=?"
	err := db.QueryExecutor(query, &temp, buyId)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "buyStore", s)
		return
	}
	fmt.Println("Tmep is ", temp)

	var playerDetails model.Player
	//get playerDetails
	query = "SELECT * FROM players WHERE player_id=?"
	err = db.QueryExecutor(query, &playerDetails, playerId)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "buyStore", s)
		return
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

	utils.SocketServerInstance.BroadcastToRoom("/", playerId, "playerDetails", playerDetails)

	response.SocketResponse(utils.PURCHASE_SUCCESS, utils.HTTP_OK, utils.SUCCESS, temp, "buyStore", s)

}
