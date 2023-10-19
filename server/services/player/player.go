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
