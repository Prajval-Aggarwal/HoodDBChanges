package car

import (
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"main/server/utils"

	socketio "github.com/googollee/go-socket.io"
)

func UpgradeCarService(s socketio.Conn, reqData map[string]interface{}) {
	playerId := s.Context().(string)
	custId := reqData["custId"].(string)
	// var upgradeCost int64
	// isUpgradable := true

	playerDetails, err := utils.GetPlayerDetails(playerId)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "carUpgrade", s)
		return
	}

	var carDetails model.PlayerCarCustomisation
	query := "SELECT * FROM player_car_customisations WHERE cust_id=?"
	err = db.QueryExecutor(query, &carDetails, custId)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "carUpgrade", s)
		return
	}
	var carClass int
	var maxLevel int
	var carClassORMultiplier float64
	var baseVaue int64
	var multiplier float64

	query = "SELECT c.class from cars c JOIN player_car_customisations pcc ON pcc.car_id=c.car_id WHERE pcc.cust_id=?"
	err = db.QueryExecutor(query, &carClass, custId)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "carUpgrade", s)
		return
	}

	switch carClass {
	case utils.D:
		maxLevel = int(utils.D_CLASS)
		carClassORMultiplier = float64(utils.D_CLASS_OR)
		baseVaue = 250
		multiplier = 1.15
	case utils.C:
		maxLevel = int(utils.C_CLASS)
		carClassORMultiplier = float64(utils.C_CLASS_OR)
		baseVaue = 500
		multiplier = 1.30

	case utils.B:
		maxLevel = int(utils.B_CLASS)
		carClassORMultiplier = float64(utils.B_CLASS_OR)
		baseVaue = 750
		multiplier = 1.50
	case utils.A:
		maxLevel = int(utils.A_CLASS)
		carClassORMultiplier = float64(utils.A_CLASS_OR)
		baseVaue = 1000
		multiplier = 1.70
	case utils.S:
		maxLevel = int(utils.S_CLASS)
		carClassORMultiplier = float64(utils.S_CLASS_OR)
		baseVaue = 1250
		multiplier = 2.00
	}

	if carDetails.CarLevel == int64(maxLevel) {
		// isUpgradable = false
		//part cannnot be upgraded further
		response.SocketResponse(utils.UPGRADE_REACHED_MAX_LEVEL, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "ack", s)
		return
	}

	//calculate cost using formula
	var cost = roundToNearestMultiple(float64(baseVaue)*(float64(1+carDetails.CarLevel)*multiplier), int(baseVaue))
	fmt.Println("Cost of the upgrade is:", cost)

	if playerDetails.Coins < int64(cost) {
		response.SocketResponse(utils.NOT_ENOUGH_COINS, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "carUpgrade", s)
		return
	}

	//upgrade car stats
	tx := db.BeginTransaction()

	//Updatting player coins
	playerDetails.Coins -= int64(cost)
	err = tx.Where(utils.PLAYER_ID, playerId).Updates(&playerDetails).Error
	if err != nil {
		tx.Rollback()
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "carUpgrade", s)
		return
	}

	carDetails.Power += utils.UPGRADE_POWER
	carDetails.Grip += int64(utils.UPGRADE_GRIP)
	carDetails.ShiftTime += utils.UPGRADE_SHIFT_TIME

	var ovr = utils.RoundFloat(CalculateOVR(carClassORMultiplier, float64(carDetails.Power), float64(carDetails.Grip), float64(carDetails.ShiftTime)), 2)

	carDetails.OVR = ovr

	err = tx.Where("cust_id", custId).Updates(&carDetails).Error
	if err != nil {
		tx.Rollback()
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "carUpgrade", s)
		return
	}

	if err := tx.Commit().Error; err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "carUpgrade", s)
		return
	}

	// var nextOvr = utils.RoundFloat(CalculateOVR(classRating.ORMultiplier, float64(playerCarStats.Power+utils.UPGRADE_POWER), float64(playerCarStats.Grip+int64(utils.UPGRADE_GRIP)), float64(playerCarStats.ShiftTime+utils.UPGRADE_SHIFT_TIME)), 2)
	// upgradeResp := &response.UpgradeResponse{
	// 	CarLevel:      carDetails.CarLevel + 1,
	// 	NextLevelCost: nextCost,
	// 	NewStats: response.Stats{
	// 		Power:      playerCarStats.Power,
	// 		Grip:       playerCarStats.Grip,
	// 		Weight:     playerCarStats.Weight,
	// 		ShiftTime:  playerCarStats.ShiftTime,
	// 		OVR:        ovr,
	// 		Durability: playerCarStats.Durability,
	// 	},
	// 	NextStats: response.Stats{
	// 		Power:      playerCarStats.Power + utils.UPGRADE_POWER,
	// 		Grip:       playerCarStats.Grip + int64(utils.UPGRADE_GRIP),
	// 		Weight:     playerCarStats.Weight,
	// 		ShiftTime:  playerCarStats.ShiftTime + utils.UPGRADE_SHIFT_TIME,
	// 		OVR:        nextOvr,
	// 		Durability: playerCarStats.Durability,
	// 	},
	// 	IsUpgradable: isUpgradable,
	// }
	// utils.SocketServerInstance.BroadcastToRoom("/", playerId, "playerDetails", playerDetails)
	// utils.SocketServerInstance.BroadcastToRoom("/", playerId, "carDetails", playerCarStats)

	response.SocketResponse(utils.UPGRADE_SUCCESS, utils.HTTP_OK, utils.SUCCESS, nil, "carUpgrade", s)

}
func CalculateOVR(classOr, power, grip, weight float64) float64 {
	x := (classOr * (0.7*float64(power) + (0.6 * float64(grip)))) - 0.02*float64(weight)
	return x
}

func roundToNearestMultiple(value float64, multiple int) float64 {
	return float64(int((value+float64(multiple/2))/float64(multiple)) * multiple)
}
