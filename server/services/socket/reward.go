package socket

import (
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"main/server/utils"
	"time"

	socketio "github.com/googollee/go-socket.io"
)

func GiveArenaPerks2(server *socketio.Server) {

	var tempRes []struct {
		PlayerId       string
		ArenaId        string
		NextRewardTime time.Time
		ArenaLevel     int64
	}
	query := `SELECT ar.player_id,ar.arena_id,a.level FROM arena_rewards  JOIN arenas a ON ar.arena_id=a.arena_id WHERE next_reward_time= CURRENT_TIMESTAMP`

	// query = "SELECT ar.*,a.level FROM arena_rewards  JOIN arenas a ON ar.arena_id= a.arena_id WHERE next_reward_time= CURRENT_TIMESTAMP"

	err := db.QueryExecutor(query, &tempRes)
	if err != nil {
		fmt.Println("Error is:", err.Error())
		return
	}

	if len(tempRes) != 0 {
		for _, temp := range tempRes {

			var nextTime time.Time

			switch int64(temp.ArenaLevel) {
			case int64(utils.EASY):
				nextTime = time.Now().Add(time.Duration(utils.EASY_PERK_MINUTES) * time.Minute)
			case int64(utils.MEDIUM):
				nextTime = time.Now().Add(time.Duration(utils.MEDIUM_PERK_MINUTES) * time.Minute)
			case int64(utils.HARD):
				nextTime = time.Now().Add(time.Duration(utils.HARD_PERK_MINUTES) * time.Minute)

			}

			var arenaPerks model.ArenaLevelPerks
			query := "SELECT * FROM arena_level_perks WHERE arena_level=?"
			err := db.QueryExecutor(query, &arenaPerks, temp.ArenaLevel)
			if err != nil {
				fmt.Println("Error is:", err.Error())
				return
			}
			var arenaRewards model.ArenaReward
			query = "SELECT * FROM arena_reward WHERE arena_id=? AND player_id=?"
			err = db.QueryExecutor(query, &arenaRewards, temp.ArenaId, temp.PlayerId)
			if err != nil {
				fmt.Println("Error in fetching from arena reward")
				return
			}

			arenaRewards.Coins += arenaPerks.Coins
			arenaRewards.Cash += arenaPerks.Cash
			arenaRewards.RepairCurrency += arenaPerks.RepairParts
			arenaRewards.RewardTime = time.Now()
			arenaRewards.NextRewardTime = nextTime

			playerDetails, _ := utils.GetPlayerDetails(temp.PlayerId)

			playerDetails.Coins += arenaPerks.Coins
			playerDetails.Cash += arenaPerks.Cash
			playerDetails.RepairCurrency += arenaPerks.RepairParts

			err = db.UpdateRecord(&playerDetails, temp.PlayerId, utils.PLAYER_ID).Error
			if err != nil {
				fmt.Println("Error is:", err.Error())
				return
			}

			utils.SocketServerInstance.BroadcastToRoom("/", playerDetails.PlayerId, "reward", arenaRewards)
		}
	} else {
		fmt.Println("No on owns the arena")
		return
	}

}

func Close(s socketio.Conn, req map[string]interface{}) {
	playerId := s.Context().(string)

	arenaId, ok := req["arenaId"].(string)
	if !ok {
		response.SocketResponse("Arena Id required", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "close", s)
		return

	}

	query := "UPDATE arena_rewards set coins=0, cash=0 , repair_currency=0 WHERE arena_id=? AND player_id=?"
	err := db.RawExecutor(query, arenaId, playerId)
	if err != nil {

		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "close", s)
		return
	}

	response.SocketResponse(utils.SUCCESS, utils.HTTP_OK, utils.SUCCESS, nil, "close", s)

}

func Open(s socketio.Conn, req map[string]interface{}) {
	playerId := s.Context().(string)

	arenaId, ok := req["arenaId"].(string)
	if !ok {
		response.SocketResponse("Arena id required", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "open", s)
		return
	}

	var arenaReward model.ArenaReward
	query := "SELECT * FROM arena_rewards WHERE arena_id=? AND player_id=?"
	err := db.QueryExecutor(query, &arenaReward, arenaId, playerId)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "open", s)
		return
	}

	response.SocketResponse(utils.DATA_FETCH_SUCCESS, utils.HTTP_OK, utils.SUCCESS, arenaReward, "open", s)

}
