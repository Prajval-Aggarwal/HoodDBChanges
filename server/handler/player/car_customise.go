package handler

import (
	"main/server/services/car"

	"github.com/gin-gonic/gin"
)

// @Summary Get the price of car parts customizations
// @Description Get the price of car parts customizations
// @Accept  json
// @Tags Car-Customize
// @Produce  json
// @Success 200 {object} response.Success
// @Router /car/customise/price [get]
func GetCustomisationPriceHandler(ctx *gin.Context) {
	car.GetCustomisePrice(ctx)
}
