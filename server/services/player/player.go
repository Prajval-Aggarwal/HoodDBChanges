package player

import (
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

func GetLevelService(ctx *gin.Context) {

	var levelResp = []response.Level{}

	var levelDetails []model.PlayerLevel
	err := db.FindAll(&levelDetails)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	//reduce tis loop
	for i := 0; i < len(levelDetails)-1; i++ {

		tempResp := response.Level{
			LastLevelXp:   levelDetails[i].XPRequired,
			NextLevelXp:   levelDetails[i+1].XPRequired,
			CurrentReward: levelDetails[i].Coins,
			RewardValue:   levelDetails[i+1].Coins,
			LevelNumber:   levelDetails[i].Level,
		}

		levelResp = append(levelResp, tempResp)
	}

	res := response.DataResponse{
		Data: levelResp,
	}
	response.ShowResponse(utils.DATA_FETCH_SUCCESS, utils.HTTP_OK, utils.SUCCESS, res, ctx)

}

func GetCarService(ctx *gin.Context, playerId string) {

	var playerCar []response.CarCustom
	query := `SELECT pcc.*,c.car_name,c.class 
	FROM player_car_customisations pcc 
	JOIN cars c ON c.car_id=pcc.car_id 
	WHERE player_id=?`
	err := db.QueryExecutor(query, &playerCar, playerId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	var resp []response.CarRes
	for _, car := range playerCar {
		temp := &response.CarRes{
			CustId:  car.CustId,
			CarId:   car.CarId,
			CarName: car.CarName,
			Rarity:  car.Class,
		}
		temp.CarCurrentData.Stats.Grip = car.Grip
		temp.CarCurrentData.Stats.Weight = car.Weight
		temp.CarCurrentData.Stats.Power = car.Power
		temp.CarCurrentData.Stats.ShiftTime = car.ShiftTime
		temp.CarCurrentData.Stats.OVR = car.OVR
		temp.CarCurrentData.Stats.Durability = car.Durability

		carCustomise, _ := utils.CustomiseMapping(car.CustId, "player_car_customisations")

		temp.CarCurrentData.Customization = *carCustomise
		temp.Status.Owned = true
		temp.Status.Purchasable = false

		resp = append(resp, *temp)

	}

	response.ShowResponse(utils.DATA_FETCH_SUCCESS, utils.HTTP_OK, utils.SUCCESS, resp, ctx)
}
