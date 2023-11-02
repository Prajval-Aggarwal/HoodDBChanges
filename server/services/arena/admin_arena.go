package arena

import (
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/request"
	"main/server/response"
	"main/server/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func AddArenaService(ctx *gin.Context, addArenaReq request.AddArenaRequest) {
	//	var newArena model.Arena

	var exists bool
	//check that no two same Arenas are on same locations
	query := "SELECT EXISTS (SELECT * FROM arenas WHERE latitude=? AND longitude=?)"
	err := db.QueryExecutor(query, &exists, addArenaReq.Latitude, addArenaReq.Longitude)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	if exists {
		response.ShowResponse(utils.ARENA_ALREADY_PRESENT, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	newArena := model.Arena{
		ArenaName:  addArenaReq.ArenaName,
		Latitude:   addArenaReq.Latitude,
		Longitude:  addArenaReq.Longitude,
		ArenaLevel: addArenaReq.ArenaLevel,
	}

	err = db.CreateRecord(&newArena)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	//add the arena to the owned arena list of the AI

	//add the win records for that arena
	err = utils.GiveArenaToAi(newArena.ArenaId, newArena.ArenaLevel)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	response.ShowResponse(utils.ARENA_ADD_SUCCESS, utils.HTTP_OK, utils.SUCCESS, newArena, ctx)
}

func DeleteArenaService(ctx *gin.Context, deleteReq request.DeletArenaReq) {
	if !db.RecordExist("arenas", deleteReq.ArenaId, "arena_id") {
		response.ShowResponse(utils.ARENA_NOT_FOUND, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}
	//validate the Arena id

	var role string
	query := "SELECT role FROM players p JOIN player_race_stats oba ON oba.player_id=p.player_id WHERE oba.arena_id=?"

	err := db.QueryExecutor(query, &role, deleteReq.ArenaId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	fmt.Println("Role is", role)
	if role != "ai" {
		response.ShowResponse("Arena is owned by players. Unable to delete the arena", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	query = "DELETE FROM arenas WHERE arena_id =?"
	err = db.RawExecutor(query, deleteReq.ArenaId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	response.ShowResponse(utils.ARENA_DELETE_SUCCESS, utils.HTTP_OK, utils.SUCCESS, nil, ctx)

}

func UpdateArenaService(ctx *gin.Context, updateReq request.UpdateArenaReq) {
	var ArenaDetails model.Arena

	//check if the Arena exists or not
	if !db.RecordExist("arenas", updateReq.ArenaId, "arena_id") {
		response.ShowResponse(utils.ARENA_NOT_FOUND, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	err := db.FindById(&ArenaDetails, updateReq.ArenaId, "arena_id")
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	//Null check on the inputs
	if updateReq.ArenaName != "" {
		ArenaDetails.ArenaName = updateReq.ArenaName
	}

	if updateReq.Latitude != 0 {
		ArenaDetails.Latitude = updateReq.Latitude
	}
	if updateReq.Longitude != 0 {
		ArenaDetails.Longitude = updateReq.Longitude
	}
	if updateReq.ArenaLevel != 0 {
		ArenaDetails.ArenaLevel = updateReq.ArenaLevel
	}

	err = db.UpdateRecord(&ArenaDetails, updateReq.ArenaId, "arena_id").Error
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	response.ShowResponse(utils.ARENA_UPDATE_SUCCESS, utils.HTTP_OK, utils.SUCCESS, ArenaDetails, ctx)

}
func GetAllArenaService(ctx *gin.Context) {
	var ArenaList = []model.Arena{}
	var dataresp response.DataResponse
	// Get the query parameters for skip and limit from the request
	skipParam := ctx.DefaultQuery("skip", "0")
	limitParam := ctx.DefaultQuery("limit", "10")

	// Convert skip and limit to integers
	skip, err := strconv.Atoi(skipParam)
	if err != nil {
		response.ShowResponse("Invalid skip value", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		response.ShowResponse("Invalid limit value", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	// Build the SQL query with skip and limit
	query := fmt.Sprintf("SELECT * FROM arenas ORDER BY created_at DESC LIMIT %d OFFSET %d", limit, skip)

	err = db.QueryExecutor(query, &ArenaList)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	var arenaResponseList []response.ArenaResponse

	for _, arena := range ArenaList {
		var arenaReward model.ArenaLevelPerks
		query := "SELECT * FROM arena_level_perks WHERE arena_level=?"
		err = db.QueryExecutor(query, &arenaReward, arena.ArenaLevel)
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}
		arenaResponse := response.ArenaResponse{
			ArenaId:    arena.ArenaId,
			ArenaName:  arena.ArenaName,
			ArenaLevel: arena.ArenaLevel,
			Longitude:  arena.Longitude,
			Latitude:   arena.Latitude,
			RewardData: arenaReward,
		}
		var res time.Duration
		switch arena.ArenaLevel {
		case int64(utils.EASY):
			// fmt.Println(time.Duration(utils.EASY_PERK_MINUTES) * time.Minute)
			res = time.Duration(utils.EASY_PERK_MINUTES) * time.Minute
			arenaResponse.NumberOfRaces = 3

		case int64(utils.MEDIUM):
			res = time.Duration(utils.MEDIUM_PERK_MINUTES) * time.Minute
			arenaResponse.NumberOfRaces = 5
		case int64(utils.HARD):
			res = time.Duration(utils.HARD_PERK_MINUTES) * time.Minute
			arenaResponse.NumberOfRaces = 7
		}

		arenaResponse.RewardTime = res.String()

		arenaResponseList = append(arenaResponseList, arenaResponse)

	}

	var totalCount int
	countQuery := "SELECT COUNT(*) FROM arenas"
	err = db.QueryExecutor(countQuery, &totalCount)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	dataresp.TotalCount = totalCount
	dataresp.Data = arenaResponseList

	response.ShowResponse(utils.DATA_FETCH_SUCCESS, utils.HTTP_OK, utils.SUCCESS, dataresp, ctx)
}

func GetArenaTypes(ctx *gin.Context) {

	var arenaTypeList = []struct {
		TypeName string `json:"label,omitempty" gorm:"unique"`
		TypeId   int    `json:"value"`
	}{
		struct {
			TypeName string "json:\"label,omitempty\" gorm:\"unique\""
			TypeId   int    "json:\"value\""
		}{
			TypeName: "Easy",
			TypeId:   int(utils.EASY),
		},
		struct {
			TypeName string "json:\"label,omitempty\" gorm:\"unique\""
			TypeId   int    "json:\"value\""
		}{
			TypeName: "Medium",
			TypeId:   int(utils.MEDIUM),
		}, struct {
			TypeName string "json:\"label,omitempty\" gorm:\"unique\""
			TypeId   int    "json:\"value\""
		}{
			TypeName: "Hard",
			TypeId:   int(utils.HARD),
		},
	}
	var dataresp response.DataResponse
	// Get the query parameters for skip and limit from the request

	dataresp.TotalCount = len(arenaTypeList)
	dataresp.Data = arenaTypeList

	response.ShowResponse(utils.DATA_FETCH_SUCCESS, utils.HTTP_OK, utils.SUCCESS, dataresp, ctx)
}
