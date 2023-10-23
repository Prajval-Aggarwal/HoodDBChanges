package car

import (
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

func GetCustomisePrice(ctx *gin.Context) {
	resp := struct {
		Color struct {
			Paint  int64 `json:"paint"`
			Livery int64 `json:"livery"`
		} `json:"color"`
		Interior int64 `json:"interior"`
		Wheels   int64 `json:"wheels"`
		Plate    int64 `json:"plate"`
	}{
		Color: struct {
			Paint  int64 `json:"paint"`
			Livery int64 `json:"livery"`
		}{
			Paint:  1000,
			Livery: 2500,
		},
		Interior: 1000,
		Wheels:   1000,
		Plate:    500,
	}
	fmt.Println("res is  ", resp)
	response.ShowResponse(utils.DATA_FETCH_SUCCESS, utils.HTTP_OK, utils.SUCCESS, resp, ctx)
}
func GetAllCarsService(ctx *gin.Context) {
	var carDetails []model.Car
	query := "SELECT * FROM cars "
	err := db.QueryExecutor(query, &carDetails)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	var finalRes struct {
		MaxStats response.Stats    `json:"maxStats,omitempty"`
		Res      []response.CarRes `json:"carData,omitempty"`
	}

	var res []response.CarRes

	for _, car := range carDetails {
		record := &response.CarRes{
			CarId:   car.CarId,
			CarName: car.CarName,
			Rarity:  car.Class,
		}
		var carStats response.Stats
		query = "SELECT power,grip,shift_time,weight,ovr,durability,nitrous_time FROM default_customisations WHERE car_id=?"
		err := db.QueryExecutor(query, &carStats, car.CarId)
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}
		var carCustomise []response.Customization
		query = "SELECT color_category,color_type,color_name,wheel_category,wheel_color_name,interior_color_name,lp_value FROM default_customisations WHERE car_id=?"
		err = db.QueryExecutor(query, &carCustomise, car.CarId)
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}
		record.Defaults.Stats = carStats
		record.Defaults.Customization = carCustomise
		record.Status.Owned = false
		record.Status.Purchasable = true
		record.Defaults.Purchase.Amount = int64(car.CurrAmount)
		switch car.CurrType {
		case "coins":
			record.Defaults.Purchase.CurrencyType = 1
		case "cash":
			record.Defaults.Purchase.CurrencyType = 2
		}

		record.Defaults.Purchase.PremiumBuy = car.PremiumBuy

		//get teh default data
		res = append(res, *record)
	}

	finalRes.MaxStats.Power = 1000
	finalRes.MaxStats.Grip = 100
	finalRes.MaxStats.ShiftTime = 10
	finalRes.MaxStats.Weight = 1000
	finalRes.MaxStats.OVR = 10
	finalRes.MaxStats.Durability = 100
	finalRes.Res = res

	response.ShowResponse(utils.DATA_FETCH_SUCCESS, utils.HTTP_OK, utils.SUCCESS, finalRes, ctx)
}
