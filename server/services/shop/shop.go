package shop

import (
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

func GetShopDetails(ctx *gin.Context) {
	var shopDetails []model.Shop

	err := db.FindAll(&shopDetails)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	response.ShowResponse(utils.DATA_FETCH_SUCCESS, utils.HTTP_OK, utils.SUCCESS, shopDetails, ctx)
}
