package handler

import (
	"main/server/services/arena"

	"github.com/gin-gonic/gin"
)

// @Summary Get arena owner
// @Description Get the details of arena owner
// @Tags Arena
// @Accept json
// @Produce json
// @Param id query string true "Id of the arena"
// @Success 200 {object} response.Success "Success"
// @Failure 400 {object} response.Success "Bad request"
// @Failure  401 {object} response.Success "Unauthorised"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /arena/owner [get]
func GetArenaOwnerHandler(ctx *gin.Context) {
	arenaId := ctx.Request.URL.Query().Get("id")
	arena.GetArenaOwnerService(ctx, arenaId)
}
