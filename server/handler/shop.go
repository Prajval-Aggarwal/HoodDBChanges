package handler

import (
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/services/shop"
	"main/server/utils"
	"math"

	"github.com/gin-gonic/gin"
)

// @Summary Get things in shop
// @Description Retrieve a list of Shop items
// @Accept  json
// @Tags Shop
// @Produce  json
// @Success 200 {object} response.Success
// @Failure 400 {object} response.Success
// @Failure 500 {object} response.Success
// @Router /get-shop [get]
func GetShopHandler(ctx *gin.Context) {
	shop.GetShopDetails(ctx)
}

func AddShopDataToDB() {

	// 1 repair part =  5 coins
	baseValue := 5

	if !utils.TableIsEmpty("shops") {

		//coins se re
		repairParts := []int64{450, 750, 1250, 1750, 2250, 3000}

		//1 for coins, 2 for cash, 3 for repairpart and 4 for real money
		for _, rp := range repairParts {
			newRecord := &model.Shop{
				PurchaseType:  utils.COINS,
				RewardType:    utils.REPAIR_PARTS,
				PurchaseValue: rp * int64(baseValue),
				RewardValue:   rp,
			}
			err := db.CreateRecord(newRecord)
			if err != nil {
				fmt.Println("Error is:", err.Error())
			}
		}

		//cash se coins buy kr rhe
		coins := []int64{4500, 7500, 10000, 50000, 75000, 125000}
		for _, coin := range coins {
			newRecord := &model.Shop{
				PurchaseType: utils.CASH,
				RewardType:   utils.COINS,
				//rouding the calculation to nearest multiple of 10
				PurchaseValue: int64((math.Round(float64(coin/15) / 10.0)) * 10.0),
				RewardValue:   coin,
			}
			err := db.CreateRecord(newRecord)
			if err != nil {
				fmt.Println("Error is:", err.Error())
			}
		}

	}
}
