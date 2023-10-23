package socket

import (
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"main/server/utils"

	socketio "github.com/googollee/go-socket.io"
)

func AmountCheckSum(s socketio.Conn, req map[string]interface{}) {
	playerId := s.Context().(string)
	id, ok := req["id"].(string)
	if !ok {
		response.SocketResponse("id is required", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "checkSum", s)
		return
	}
	name, ok := req["name"].(string)
	if !ok {
		response.SocketResponse("Name is required", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "checkSum", s)
		return
	}
	purchaseType, ok := req["purchaseType"].(float64)
	if !ok {
		response.SocketResponse("Purchase type is required", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "carBuy", s)
		return
	}

	var buyType string
	switch purchaseType {
	case utils.COINS:
		buyType = "coins"
	case utils.CASH:
		buyType = "cash"
	}
	//shift it to getPlayerDetails function as it is commonly used

	playerDetails, err := utils.GetPlayerDetails(playerId)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "checkSum", s)
		return
	}

	switch name {
	case "car":
		{
			var details model.Car
			query := "SELECT * FROM cars WHERE car_id=?"
			err := db.QueryExecutor(query, &details, id)
			if err != nil {
				response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "checkSum", s)
				return
			}

			if buyType == details.CurrType {
				if details.CurrType == "coins" {
					if details.CurrAmount > float64(playerDetails.Coins) {
						response.SocketResponse(utils.NOT_ENOUGH_COINS, utils.HTTP_BAD_REQUEST, utils.FAILURE, req, "checkSum", s)
						return
					}
				} else if details.CurrType == "cash" {
					if details.CurrAmount > float64(playerDetails.Cash) {
						response.SocketResponse(utils.NOT_ENOUGH_CASH, utils.HTTP_BAD_REQUEST, utils.FAILURE, req, "checkSum", s)
						return
					}
				}
			} else {
				if details.CurrAmount > float64(playerDetails.Cash) {
					response.SocketResponse(utils.NOT_ENOUGH_CASH, utils.HTTP_BAD_REQUEST, utils.FAILURE, req, "checkSum", s)
					return
				}
			}
			if details.CurrType == "coins" {
				if details.CurrAmount > float64(playerDetails.Coins) {
					response.SocketResponse(utils.NOT_ENOUGH_COINS, utils.HTTP_BAD_REQUEST, utils.FAILURE, req, "checkSum", s)
					return
				}
			} else if details.CurrType == "cash" {
				if details.CurrAmount > float64(playerDetails.Cash) {
					response.SocketResponse(utils.NOT_ENOUGH_CASH, utils.HTTP_BAD_REQUEST, utils.FAILURE, req, "checkSum", s)
					return
				}
			}

		}
	case "garage":
		{
			var details model.Garage
			query := "SELECT * FROM garages WHERE car_id=?"
			err := db.QueryExecutor(query, &details, id)
			if err != nil {
				response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "checkSum", s)
				return
			}

			if float64(details.CoinsRequired) > float64(playerDetails.Coins) {
				response.SocketResponse(utils.NOT_ENOUGH_COINS, utils.HTTP_BAD_REQUEST, utils.FAILURE, req, "checkSum", s)
				return
			}

		}
	case "shop":
		{
			var details model.Shop
			query := "SELECT * FROM shops WHERE car_id=?"
			err := db.QueryExecutor(query, &details, id)
			if err != nil {
				response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "checkSum", s)
				return
			}

			switch details.PurchaseType {
			case utils.COINS:
				{
					if details.PurchaseValue > playerDetails.Coins {
						response.SocketResponse(utils.NOT_ENOUGH_COINS, utils.HTTP_BAD_REQUEST, utils.FAILURE, req, "checkSum", s)
					}
				}
			case utils.CASH:
				{
					if details.PurchaseValue > playerDetails.Cash {
						response.SocketResponse(utils.NOT_ENOUGH_CASH, utils.HTTP_BAD_REQUEST, utils.FAILURE, req, "checkSum", s)
					}
				}
			case utils.REPAIR_PARTS:
				{
					if details.PurchaseValue > playerDetails.RepairCurrency {
						response.SocketResponse(utils.NOT_ENOUGH_REPAIR_PARTS, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "checkSum", s)
					}
				}
			}

		}

	}

	response.SocketResponse(utils.SUCCESS, utils.HTTP_OK, utils.SUCCESS, req, "checkSum", s)
}
