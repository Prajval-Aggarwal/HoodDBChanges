package arena

import (
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/request"
	"main/server/response"
	"main/server/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func GetArenaSlotDetailsService(ctx *gin.Context, playerId string, arenaId string) {
	var arenaSlotData response.ArenaSlotResponse

	//check if the areana is owned by the player or not

	var arenaRewardDetails model.ArenaLevelPerks
	query := "select ap.* from arena_level_perks ap JOIN arenas a ON a.arena_level=ap.arena_level WHERE a.arena_id=?;"
	err := db.QueryExecutor(query, &arenaRewardDetails, arenaId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	arenaSlotData.RewardData = arenaRewardDetails
	if arenaRewardDetails.ArenaLevel == int64(utils.EASY) {
		arenaSlotData.TotalSlots = int(utils.EASY_ARENA_SLOT)
	} else if arenaRewardDetails.ArenaLevel == int64(utils.MEDIUM) {
		arenaSlotData.TotalSlots = int(utils.MEDIUM_ARENA_SLOT)
	} else if arenaRewardDetails.ArenaLevel == int64(utils.HARD) {
		arenaSlotData.TotalSlots = int(utils.HARD_ARENA_SLOT)
	}

	//get the time left for slots to fill
	var arenaWinTime time.Time
	query = "SELECT win_time FROM player_race_stats WHERE arena_id=? and player_id=? order by updated_at DESC LIMIT 1"
	err = db.QueryExecutor(query, &arenaWinTime, arenaId, playerId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	// temp := time.Until(arenaWinTime.Add(24 * time.Hour))
	temp := time.Until(arenaWinTime.Add(3 * time.Minute))

	// temp := arenaWinTime.Add(24 * time.Hour).Sub(time.Now())
	arenaSlotData.ArenaWinTime = temp.String()

	//get the next arena perk time
	var arenaPerkTime time.Time
	query = "SELECT next_reward_time FROM arena_rewards WHERE arena_id=? AND player_id=?"
	err = db.QueryExecutor(query, &arenaPerkTime, arenaId, playerId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	temp1 := arenaPerkTime.Sub(time.Now())
	fmt.Println("arena perk time is", temp1)
	arenaSlotData.ArenaPerkTime = temp1.String()

	var res []response.CarRes

	var carStruct2 []response.CarCustom
	query = `SELECT pcc.*,c.class,c.car_name
        FROM player_car_customisations pcc
        JOIN cars c ON c.car_id=pcc.car_id 
        JOIN arena_cars arr ON arr.cust_id=pcc.cust_id
        WHERE arr.arena_id=? AND arr.player_id=?;`

	// Execute the query and store the result in 'carStruct2'
	err = db.QueryExecutor(query, &carStruct2, arenaId, playerId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	for _, details := range carStruct2 {
		record := &response.CarRes{
			CustId:  details.CustId,
			CarId:   details.CarId,
			CarName: details.CarName,
			Rarity:  details.Class,
		}

		carCustomise, _ := utils.CustomiseMapping(details.CustId, "player_car_customisations")
		record.CarCurrentData.Customization = *carCustomise
		record.CarCurrentData.Stats.Power = details.Power
		record.CarCurrentData.Stats.Grip = details.Grip
		record.CarCurrentData.Stats.Weight = details.Weight
		record.CarCurrentData.Stats.ShiftTime = details.ShiftTime
		record.CarCurrentData.Stats.OVR = details.OVR
		record.CarCurrentData.Stats.Durability = details.Durability
		record.Status.Owned = true

		res = append(res, *record)
	}
	arenaSlotData.CarDetails = res

	response.ShowResponse(utils.DATA_FETCH_SUCCESS, utils.HTTP_OK, utils.SUCCESS, arenaSlotData, ctx)

}

func AddCarToSlotService(ctx *gin.Context, addCarReq request.AddCarArenaRequest, playerId string) {

	// Check if the car is bought by the player
	query := "SELECT EXISTS(SELECT * FROM owned_cars WHERE player_id = ? AND cust_id = ?)"
	if !utils.IsExisting(query, playerId, addCarReq.CustId) {
		response.ShowResponse(utils.CAR_NOT_OWNED, utils.HTTP_NOT_FOUND, utils.FAILURE, nil, ctx)
		return
	}

	// Check if the player owns the arena
	query = "SELECT EXISTS(SELECT * FROM player_race_stats WHERE player_id = ? AND arena_id = ?)"
	if !utils.IsExisting(query, playerId, addCarReq.ArenaId) {
		response.ShowResponse(utils.ARENA_NOT_OWNED, utils.HTTP_NOT_FOUND, utils.FAILURE, nil, ctx)
		return
	}

	// Check that it should not add more cars than required slots for the arena
	var arenaDetails model.Arena
	err := db.FindById(&arenaDetails, addCarReq.ArenaId, "arena_id")
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	//check that if the car is already alloted to another arena or not
	query = "SELECT EXISTS (SELECT * FROM arena_cars WHERE player_id = ? AND cust_id=?)"
	if utils.IsExisting(query, playerId, addCarReq.CustId) {
		response.ShowResponse(utils.CAR_ALREADY_ALLOTTED, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	var carCount int64
	query = "SELECT COUNT(*) FROM arena_cars WHERE player_id = ? AND arena_id = ?"
	err = db.QueryExecutor(query, &carCount, playerId, addCarReq.ArenaId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	// Check the slot limit for the arena level and ensure it's not exceeded
	var maxSlots int64
	switch arenaDetails.ArenaLevel {
	case int64(utils.EASY):
		maxSlots = utils.EASY_ARENA_SLOT
	case int64(utils.MEDIUM):
		maxSlots = utils.MEDIUM_ARENA_SLOT
	case int64(utils.HARD):
		maxSlots = utils.HARD_ARENA_SLOT
	default:
		response.ShowResponse("Invalid arena level", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	if carCount == maxSlots {
		response.ShowResponse(utils.NO_CARS_ADDED, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	// Create a record in the car_slots table
	carSlot := model.ArenaCars{
		PlayerId: playerId,
		ArenaId:  addCarReq.ArenaId,
		CustId:   addCarReq.CustId,
	}

	err = db.CreateRecord(&carSlot)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	response.ShowResponse(utils.CAR_ADDED_SUCCESS, utils.HTTP_OK, utils.SUCCESS, carSlot, ctx)
}

func ReplaceCarService(ctx *gin.Context, replaceReq request.ReplaceReq, playerId string) {
	// Check if the car is bought by the player and owned by the player
	query := "SELECT EXISTS(SELECT * FROM owned_cars WHERE player_id = ? AND cust_id = ?)"
	if !utils.IsExisting(query, playerId, replaceReq.NewCustId) {
		response.ShowResponse(utils.CAR_NOT_OWNED, utils.HTTP_NOT_FOUND, utils.FAILURE, nil, ctx)
		return
	}

	// Check if the player owns the arena
	query = "SELECT EXISTS(SELECT * FROM player_race_stats WHERE player_id = ? AND arena_id = ?)"
	if !utils.IsExisting(query, playerId, replaceReq.ArenaId) {
		response.ShowResponse(utils.ARENA_NOT_OWNED, utils.HTTP_NOT_FOUND, utils.FAILURE, nil, ctx)
		return
	}

	//check that if the car is already alloted to another arena or not
	query = "SELECT EXISTS (SELECT * FROM arena_cars WHERE player_id = ? AND cust_id=?)"
	if utils.IsExisting(query, playerId, replaceReq.NewCustId) {
		response.ShowResponse(utils.CAR_ALREADY_ALLOTTED, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	// Replace the car in the slot
	query = "UPDATE arena_cars SET cust_id = ? WHERE player_id = ? AND arena_id = ? AND cust_id=?"
	err := db.RawExecutor(query, replaceReq.NewCustId, playerId, replaceReq.ArenaId, replaceReq.ExistingCustId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	response.ShowResponse(utils.CAR_REPLACED_SUCCESS, utils.HTTP_OK, utils.SUCCESS, replaceReq, ctx)
}
