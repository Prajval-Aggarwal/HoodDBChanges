package handler

import (
	"main/server/request"
	"main/server/response"
	"main/server/services/arena"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

// @Summary End Challenge
// @Description Ends the current challenge and saves the data
// @Tags Arena
// @Accept json
// @Produce json
// @Param Authorization header string true "Player Access token"
// @Param challengereq body request.EndChallengeReq true "End Challenge Request"
// @Success 200 {object} response.Success "Success"
// @Failure 400 {object} response.Success "Bad request"
// @Failure  401 {object} response.Success "Unauthorised"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /arena/end [post]
func EndChallengeHandler(ctx *gin.Context) {
	playerId, exists := ctx.Get(utils.PLAYERID)
	if !exists {

		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}
	var endChallReq request.EndChallengeReq
	err := utils.RequestDecoding(ctx, &endChallReq)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	err = endChallReq.Validate()
	if err != nil {

		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}
	arena.EndChallengeService(ctx, endChallReq, playerId.(string))
}

// @Summary Enter the arena
// @Description Enter arena
// @Tags Arena
// @Accept json
// @Produce json
// @Param Authorization header string true "Player Access token"
// @Param addCarReq body request.GetArenaReq true "Id of the arena"
// @Success 200 {object} response.Success "Success"
// @Failure 400 {object} response.Success "Bad request"
// @Failure  401 {object} response.Success "Unauthorised"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /arena/enter [post]
func EnterArenaHandler(ctx *gin.Context) {
	playerId, exists := ctx.Get(utils.PLAYERID)
	if !exists {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}
	var enterReq request.GetArenaReq
	err := utils.RequestDecoding(ctx, &enterReq)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	err = enterReq.Validate()
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	arena.EnterArenaService(ctx, enterReq, playerId.(string))
}

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

// AddCarToSlotHandler adds a car to the player's slot in a specific arena.
// @Summary Add a car to arena
// @Description Add a car to the player's slot in a specific arena
// @Tags Arena
// @Accept json
// @Produce json
// @Param Authorization header string true "Player Access token"
// @Param addCarReq body request.AddCarArenaRequest true "Add car to slot request payload"
// @Success 200 {object} response.Success "Car added to slot successfully"
// @Failure 400 {object} response.Success "Bad request. Invalid payload"
// @Failure 404 {object} response.Success "Car or player not found"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /arena/add-car [post]
func AddCarToSlotHandler(ctx *gin.Context) {
	playerId, exists := ctx.Get(utils.PLAYERID)
	if !exists {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}
	var addCarReq request.AddCarArenaRequest
	err := utils.RequestDecoding(ctx, &addCarReq)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	arena.AddCarToSlotService(ctx, addCarReq, playerId.(string))

}

// @Summary Replace Car
// @Description Add or replaces the car in the arena car slot
// @Tags Arena
// @Accept json
// @Produce json
// @Param Authorization header string true "Player Access token"
// @Param challengereq body request.ReplaceReq true "Replace car Request"
// @Success 200 {object} response.Success "Success"
// @Failure 400 {object} response.Success "Bad request"
// @Failure  401 {object} response.Success "Unauthorised"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /arena/replace-car [put]
func ReplaceCarHandler(ctx *gin.Context) {
	playerId, exists := ctx.Get(utils.PLAYERID)
	if !exists {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}
	var addCarReq request.ReplaceReq
	err := utils.RequestDecoding(ctx, &addCarReq)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	arena.ReplaceCarService(ctx, addCarReq, playerId.(string))
}

// ArenaCarHandler Get the list of cars that can enter the arena
// @Summary List of cars that can play
// @Description Get the list of cars that can enter the arena
// @Tags Arena
// @Accept json
// @Produce json
// @Param Authorization header string true "Player Access token"
// @Success 200 {object} response.Success "Car added to slot successfully"
// @Failure 400 {object} response.Success "Bad request. Invalid payload"
// @Failure 404 {object} response.Success "Car or player not found"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /arena/cars [get]
func ArenaCarHandler(ctx *gin.Context) {
	playerId, exists := ctx.Get(utils.PLAYERID)
	if !exists {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}
	arena.ArenaCarService(ctx, playerId.(string))
}

// @Summary Get arena slots details
// @Description Get the details of the cars kept in arena
// @Tags Arena
// @Accept json
// @Produce json
// @Param Authorization header string true "Player Access token"
// @Param id query string true "Id of the arena"
// @Success 200 {object} response.Success "Success"
// @Failure 400 {object} response.Success "Bad request"
// @Failure  401 {object} response.Success "Unauthorised"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /arena/slots/get [get]
func GetArenaSlotDetails(ctx *gin.Context) {
	playerId, exists := ctx.Get(utils.PLAYERID)
	if !exists {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}
	arenaId := ctx.Request.URL.Query().Get("id")
	arena.GetArenaSlotDetailsService(ctx, playerId.(string), arenaId)
}
