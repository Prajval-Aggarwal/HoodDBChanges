package car

import (
	"fmt"
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
