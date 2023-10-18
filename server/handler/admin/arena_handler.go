package handler

import (
	"fmt"
	"main/server/request"
	"main/server/response"
	"main/server/services/arena"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

// @Summary Add a new arena
// @Description Add a new arena to the system
// @Tags Arena
// @Accept json
// @Produce json
// @Param Authorization header string true "Admin Access token"
// @Param garageReq body request.AddArenaRequest true "Arena request payload"
// @Success 200 {object} response.Success "Arena added successful"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /admin/arena [post]
func AddArenaHandler(ctx *gin.Context) {
	role, exists := ctx.Get("role")
	fmt.Println("player id is", role)
	if !exists {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}

	if role != "admin" {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_FORBIDDEN, utils.FAILURE, nil, ctx)
		return
	}
	var addArenaReq request.AddArenaRequest
	err := utils.RequestDecoding(ctx, &addArenaReq)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	err = addArenaReq.Validate()
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	arena.AddArenaService(ctx, addArenaReq)
}

// DeleteArenaHandler deletes a arena with the given ID.
// @Summary Delete a Arena
// @Description Delete a Arena by its ID
// @Tags Arena
// @Accept json
// @Produce json
// @Param Authorization header string true "Admin Access token"
// @Param ArenaReq body request.DeletArenaReq true "Arena request payload"
// @Success 200 {object} response.Success "Arena deleted successful"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 404 {string} response.Success "Arena not found"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /admin/arena [delete]
func DeleteArenaHandler(ctx *gin.Context) {
	role, exists := ctx.Get("role")
	fmt.Println("player id is", role)
	if !exists {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}

	if role != "admin" {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_FORBIDDEN, utils.FAILURE, nil, ctx)
		return
	}
	var deleteReq request.DeletArenaReq
	err := utils.RequestDecoding(ctx, &deleteReq)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	err = deleteReq.Validate()
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	arena.DeleteArenaService(ctx, deleteReq)
}

// UpdateArenaHandler updates a Arena with the given ID.
// @Summary Update a Arena
// @Description Update a Arena by its ID
// @Tags Arena
// @Accept json
// @Produce json
// @Param Authorization header string true "Admin Access token"
// @Param updateReq body request.UpdateArenaReq true "Update request payload"
// @Success 200 {object} response.Success "Arena updated successful"
// @Failure 400 {object} response.Success "Bad request"
// @Failure 404 {string} response.Success "Arena not found"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /admin/arena [put]
func UpdateArenaHandler(ctx *gin.Context) {
	role, exists := ctx.Get("role")
	fmt.Println("player id is", role)
	if !exists {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		return
	}

	if role != "admin" {
		response.ShowResponse(utils.UNAUTHORIZED, utils.HTTP_FORBIDDEN, utils.FAILURE, nil, ctx)
		return
	}
	var updateReq request.UpdateArenaReq
	err := utils.RequestDecoding(ctx, &updateReq)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	err = updateReq.Validate()
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	arena.UpdateArenaService(ctx, updateReq)
}

// GetArenaListHandler retrieves the list of all Arenas.
//
// @Summary Get All Arena List
// @Description Retrieve the list of all Arenas
// @Tags Arena
// @Accept json
// @Produce json
// @Param skip query integer false "Number of records to skip (default is 0)"
// @Param limit query integer false "Maximum number of records to fetch (default is 10)"
// @Success 200 {object} response.Success "Arena list fetched successfully"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /arena/get [get]
func GetArenaListHandler(ctx *gin.Context) {
	arena.GetAllArenaService(ctx)
}

// GetArenaTypeHandler retrieves the list of all garages.
//
// @Summary Get All Arena type List
// @Description Retrieve the list of all arena types
// @Tags Arena
// @Accept json
// @Produce json
// @Success 200 {object} response.Success "Arena type list fetched successfully"
// @Failure 500 {object} response.Success "Internal server error"
// @Router /arena/types [get]
func GetArenaTypeHandler(ctx *gin.Context) {
	arena.GetArenaTypes(ctx)
}
