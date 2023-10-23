package handler

import (
	"main/server/services/car"

	"github.com/gin-gonic/gin"
)

// GetAllCarsServiceretrieves the list of all car.
//
// @Summary Get All Cars List
// @Description Retrieve the list of all car
// @Tags Car
// @Accept json
// @Produce json
// @Success 200 {object} response.Success "Cars list fetched successfully"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /car/get-all [get]
func GetAllCarsHandler(ctx *gin.Context) {
	car.GetAllCarsService(ctx)
}
