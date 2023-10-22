package arena

import (
	"errors"
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/request"
	"main/server/response"
	"main/server/utils"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
)

type carDetails struct {
	CustId string
	OVR    float64
}

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

	var AIId string
	query = "SELECT player_id FROM players WHERE role='ai' order by RANDOM() LIMIT 1;"
	err = db.QueryExecutor(query, &AIId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	// aiOwnedArena := model.OwnedBattleArenas{
	// 	ArenaId:  newArena.ArenaId,
	// 	PlayerId: AIId,
	// }

	aiOwnedArena := model.PlayerRaceStats{
		PlayerId: AIId,
		ArenaId:  &newArena.ArenaId,
	}
	if newArena.ArenaLevel == int64(utils.EASY) {
		aiOwnedArena.WinStreak = utils.EASY_ARENA_SERIES
	} else if newArena.ArenaLevel == int64(utils.MEDIUM) {
		aiOwnedArena.WinStreak = utils.MEDIUM_ARENA_SERIES
	} else if newArena.ArenaLevel == int64(utils.HARD) {
		aiOwnedArena.WinStreak = utils.HARD_ARENA_SERIES
	}

	err = db.CreateRecord(&aiOwnedArena)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	//add the win records for that arena
	var randomTimeSlice []string
	var carSlice []carDetails

	if addArenaReq.ArenaLevel == int64(utils.EASY) {
		randomTimeSlice = utils.GenerateRandomTime(int(utils.EASY_ARENA_SLOT), 22.0, 25.0)
		carSlice, err = GiveRandomCar(aiOwnedArena.PlayerId, newArena.ArenaId, 1, 2, int(utils.EASY_ARENA_SLOT))
		if err != nil {
			return
		}
	} else if addArenaReq.ArenaLevel == int64(utils.MEDIUM) {
		randomTimeSlice = utils.GenerateRandomTime(int(utils.MEDIUM_ARENA_SLOT), 22.0, 25.0)
		carSlice, err = GiveRandomCar(aiOwnedArena.PlayerId, newArena.ArenaId, 2, 4, int(utils.MEDIUM_ARENA_SLOT))
		if err != nil {
			return
		}

	} else if addArenaReq.ArenaLevel == int64(utils.HARD) {
		randomTimeSlice = utils.GenerateRandomTime(int(utils.HARD_ARENA_SLOT), 22.0, 25.0)
		carSlice, err = GiveRandomCar(aiOwnedArena.PlayerId, newArena.ArenaId, 4, 5, int(utils.HARD_ARENA_SLOT))
		if err != nil {
			return
		}
	}

	for i, val := range randomTimeSlice {
		newRecord := model.ArenaRaceRecord{
			PlayerId: aiOwnedArena.PlayerId,
			ArenaId:  newArena.ArenaId,
			TimeWin:  fmt.Sprintf("%v", val),
			CustId:   carSlice[i].CustId,
			Result:   "win",
		}

		err := db.CreateRecord(&newRecord)
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}
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

	var totalCount int
	countQuery := "SELECT COUNT(*) FROM arenas"
	err = db.QueryExecutor(countQuery, &totalCount)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	dataresp.TotalCount = totalCount
	dataresp.Data = ArenaList

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

func GiveRandomCar(playerId string, arenaId string, min int64, max int64, slots int) ([]carDetails, error) {

	var carSlice []carDetails
	for i := 0; i < slots; i++ {
		var carId string
		query := ` SELECT car_id FROM cars
					WHERE class >= ? AND class <= ?
					ORDER BY RANDOM() LIMIT 1 ; `
		err := db.QueryExecutor(query, &carId, min, max)
		if err != nil {
			return nil, errors.New("error in selecting the random car from the db for ai")
		}
		var carDefaults model.DefaultCustomisation
		query = "SELECT * FROM default_customisations WHERE car_id=? "
		err = db.QueryExecutor(query, &carDefaults, carId)
		if err != nil {

			return nil, err
		}

		playerCarCustomisations := model.PlayerCarCustomisation{
			PlayerId:          playerId,
			CarId:             carId,
			CarLevel:          1,
			Power:             carDefaults.Power,
			Grip:              carDefaults.Grip,
			ShiftTime:         carDefaults.ShiftTime,
			Weight:            carDefaults.Weight,
			OVR:               carDefaults.OVR,
			Durability:        carDefaults.Durability,
			NitrousTime:       carDefaults.NitrousTime,
			ColorCategory:     carDefaults.ColorCategory,
			ColorType:         carDefaults.ColorType,
			ColorName:         carDefaults.ColorName,
			WheelCategory:     carDefaults.WheelCategory,
			WheelColorName:    carDefaults.WheelColorName,
			InteriorColorName: carDefaults.InteriorColorName,
			LPValue:           carDefaults.LPValue,
		}

		err = db.CreateRecord(&playerCarCustomisations)
		if err != nil {
			return nil, err
		}

		newCarRecord := model.OwnedCars{
			PlayerId: playerId,
			CustId:   playerCarCustomisations.CustId,
			Selected: true,
		}

		err = db.CreateRecord(&newCarRecord)
		if err != nil {
			return nil, err
		}
		//Get the customisation id

		carSlice = append(carSlice, carDetails{
			CustId: playerCarCustomisations.CustId,
			OVR:    carDefaults.OVR,
		})
	}

	fmt.Println("Car slice before sorting is:", carSlice)

	sort.SliceStable(carSlice, func(i, j int) bool {
		return carSlice[i].OVR > carSlice[j].OVR
	})

	fmt.Println("Car slice is after sorting:", carSlice)

	return carSlice, nil
}
