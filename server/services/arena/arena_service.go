package arena

import (
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"main/server/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetArenaOwnerService(ctx *gin.Context, arenaId string) {

	if !db.RecordExist("arenas", arenaId, "arena_id") {
		response.ShowResponse(utils.ARENA_NOT_FOUND, utils.HTTP_NOT_FOUND, utils.FAILURE, nil, ctx)
		return
	}

	// var owner model.OwnedBattleArenas
	// query := "SELECT * FROM owned_battle_arenas WHERE arena_id=? ORDER BY updated_at DESC LIMIT 1"

	// err := db.QueryExecutor(query, &owner, arenaId)
	// if err != nil {
	// 	response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
	// 	return
	// }

	var owner model.PlayerRaceStats
	query := "SELECT * FROM player_race_stats WHERE arena_id=? and win_streak>lose_streak ORDER BY updated_at DESC LIMIT 1"
	err := db.QueryExecutor(query, &owner, arenaId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}
	var playerDetails model.Player
	query = "SELECT * FROM players WHERE player_id=?"

	err = db.QueryExecutor(query, &playerDetails, owner.PlayerId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	var arenaRecord []string

	query = "SELECT time_win from arena_race_records WHERE arena_id=? AND player_id=? ORDER BY created_at"

	err = db.QueryExecutor(query, &arenaRecord, arenaId, playerDetails.PlayerId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	var carStruct2 []struct {
		CustId            string  `json:"custId"  gorm:"unique;default:uuid_generate_v4();primaryKey,omitempty"`
		CarId             string  `json:"carId,omitempty"`
		Power             int64   `json:"power,omitempty"`
		Grip              int64   `json:"grip,omitempty"`
		ShiftTime         float64 `json:"shiftTime,omitempty"`
		Weight            int64   `json:"weight,omitempty"`
		OVR               float64 `json:"or,omitempty"` //overall rating of the car
		Durability        int64   `json:"Durability,omitempty"`
		NitrousTime       float64 `json:"nitrousTime,omitempty"` //increased when nitrous is upgraded
		ColorCategory     string  `json:"colorCategory,omitempty"`
		ColorType         string  `json:"colorType,omitempty"`
		ColorName         string  `json:"colorName,omitempty"`
		WheelCategory     string  `json:"wheelCategory,omitempty"`
		WheelColorName    string  `json:"wheelColorName,omitempty"`
		InteriorColorName string  `json:"interiorColorName,omitempty"`
		LPValue           string  `json:"lp_value,omitempty"`
	}
	var resp struct {
		PlayerId     string `json:"playerId"`
		PlayerName   string `json:"playerName"`
		ArenaId      string `json:"arenaId"`
		ArenaRecords struct {
			WinTimes []struct {
				Second       int `json:"seconds"`
				MilliSeconds int `json:"milliSeconds"`
				MicroSecond  int `json:"microSeconds"`
			} `json:"winTimes"`
		} `json:"arenaRecords"`
		Cars []struct {
			CustId            string  `json:"custId"  gorm:"unique;default:uuid_generate_v4();primaryKey,omitempty"`
			CarId             string  `json:"carId,omitempty"`
			Power             int64   `json:"power,omitempty"`
			Grip              int64   `json:"grip,omitempty"`
			ShiftTime         float64 `json:"shiftTime,omitempty"`
			Weight            int64   `json:"weight,omitempty"`
			OVR               float64 `json:"or,omitempty"` //overall rating of the car
			Durability        int64   `json:"Durability,omitempty"`
			NitrousTime       float64 `json:"nitrousTime,omitempty"` //increased when nitrous is upgraded
			ColorCategory     string  `json:"colorCategory,omitempty"`
			ColorType         string  `json:"colorType,omitempty"`
			ColorName         string  `json:"colorName,omitempty"`
			WheelCategory     string  `json:"wheelCategory,omitempty"`
			WheelColorName    string  `json:"wheelColorName,omitempty"`
			InteriorColorName string  `json:"interiorColorName,omitempty"`
			LPValue           string  `json:"lp_value,omitempty"`
		} `json:"cars"`
	}

	resp.PlayerId = playerDetails.PlayerId
	resp.PlayerName = playerDetails.PlayerName
	resp.ArenaId = arenaId

	for _, ts := range arenaRecord {

		parts := strings.Split(ts, ":")

		subParts := strings.Split(parts[2], ".")

		seconds, _ := strconv.Atoi(subParts[0])
		milliSecond, _ := strconv.Atoi(subParts[1][:2])

		microSecond, _ := strconv.Atoi(subParts[1][2:])

		resp.ArenaRecords.WinTimes = append(resp.ArenaRecords.WinTimes, struct {
			Second       int `json:"seconds"`
			MilliSeconds int `json:"milliSeconds"`
			MicroSecond  int `json:"microSeconds"`
		}{
			Second:       seconds,
			MilliSeconds: milliSecond,
			MicroSecond:  microSecond,
		})
	}

	query = `SELECT pcc.*
		FROM player_car_customisations pcc
		JOIN arena_race_records arr ON arr.cust_id=pcc.cust_id
		WHERE arr.arena_id=? AND arr.player_id=?;	`

	err = db.QueryExecutor(query, &carStruct2, arenaId, playerDetails.PlayerId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	resp.Cars = append(resp.Cars, carStruct2...)

	response.ShowResponse(utils.DATA_FETCH_SUCCESS, utils.HTTP_OK, utils.SUCCESS, resp, ctx)

}
