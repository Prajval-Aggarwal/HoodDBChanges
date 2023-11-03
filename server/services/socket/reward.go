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

type PerkResposne struct {
	ArenaId        string `json:"arenaId"`
	PlayerId       string `json:"playerId"`
	Coins          int64  `json:"coins"`
	Cash           int64  `json:"cash"`
	RepairCurrency int64  `json:"repairPart"`
	NextRewardTime string `json:"nextRewardTime"`
}

func GiveArenaPerks2(server *socketio.Server, tm time.Time) {

	var tempRes []struct {
		Id             string
		PlayerId       string
		ArenaId        string
		NextRewardTime time.Time
		ArenaLevel     int64
	}
	// query := `SELECT ar.player_id,ar.arena_id,a.arena_level FROM arena_rewards ar JOIN arenas a ON ar.arena_id=a.arena_id WHERE next_reward_time= CURRENT_TIMESTAMP`

	query := `
	SELECT ar.*, a.arena_level
	FROM arena_rewards ar
	JOIN arenas a ON ar.arena_id = a.arena_id
	WHERE
	date_trunc('minute', ar.next_reward_time) = date_trunc('minute', CURRENT_TIMESTAMP);
	`

	// query := `SELECT ar.*, a.arena_level
	// FROM arena_rewards ar
	// JOIN arenas a ON ar.arena_id = a.arena_id
	// WHERE
	// ar.next_reward_time - ? < 0 `

	// query = "SELECT ar.*,a.level FROM arena_rewards  JOIN arenas a ON ar.arena_id= a.arena_id WHERE next_reward_time= CURRENT_TIMESTAMP"

	err := db.QueryExecutor(query, &tempRes)
	if err != nil {
		fmt.Println("Error is:", err.Error())
		return
	}

	fmt.Printf("Data is%+v", tempRes)

	if len(tempRes) != 0 {
		for _, temp := range tempRes {

			// if tm.Sub(temp.NextRewardTime) {

			// }
			currentTime := time.Now()
			newTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute(), 0, 0, time.Local)
			var nextTime time.Time

			switch int64(temp.ArenaLevel) {
			case int64(utils.EASY):
				nextTime = newTime.Add(time.Duration(utils.EASY_PERK_MINUTES) * time.Minute)
			case int64(utils.MEDIUM):
				nextTime = newTime.Add(time.Duration(utils.MEDIUM_PERK_MINUTES) * time.Minute)
			case int64(utils.HARD):
				nextTime = newTime.Add(time.Duration(utils.HARD_PERK_MINUTES) * time.Minute)

			}

			var arenaPerks model.ArenaLevelPerks
			query := "SELECT * FROM arena_level_perks WHERE arena_level=?"
			err := db.QueryExecutor(query, &arenaPerks, temp.ArenaLevel)
			if err != nil {
				fmt.Println("Error is:", err.Error())
				return
			}
			var arenaRewards model.ArenaReward
			query = "SELECT * FROM arena_rewards WHERE arena_id=? AND player_id=?"
			err = db.QueryExecutor(query, &arenaRewards, temp.ArenaId, temp.PlayerId)
			if err != nil {
				fmt.Println("Error in fetching from arena reward")
				return
			}

			arenaRewards.Coins += arenaPerks.Coins
			arenaRewards.Cash += arenaPerks.Cash
			arenaRewards.RepairCurrency += arenaPerks.RepairParts
			arenaRewards.RewardTime = newTime
			arenaRewards.NextRewardTime = nextTime

			//update arena reards details
			tx := db.BeginTransaction()

			playerDetails, _ := utils.GetPlayerDetails(temp.PlayerId)

			playerDetails.Coins += arenaPerks.Coins
			playerDetails.Cash += arenaPerks.Cash
			playerDetails.RepairCurrency += arenaPerks.RepairParts

			err = tx.Where("player_id", temp.PlayerId).Updates(&playerDetails).Error
			if err != nil {
				tx.Rollback()
				fmt.Println("Error is:", err.Error())
				return
			}

			err = tx.Where("id", temp.Id).Updates(&arenaRewards).Error
			if err != nil {
				tx.Rollback()
				fmt.Println("Error is:", err.Error())
				return
			}

			res := response.Success{
				Status:  utils.SUCCESS,
				Code:    utils.HTTP_OK,
				Message: utils.SUCCESS,
				Data: PerkResposne{
					ArenaId:        arenaRewards.ArenaId,
					PlayerId:       arenaRewards.PlayerId,
					Coins:          arenaRewards.Coins,
					Cash:           arenaRewards.Cash,
					RepairCurrency: arenaRewards.RepairCurrency,
					NextRewardTime: arenaRewards.NextRewardTime.Sub(time.Now()).String(),
				},
			}

			if !utils.SocketServerInstance.BroadcastToRoom("/", playerDetails.PlayerId, "reward", res) {
				tx.Rollback()
				fmt.Println("advajsdvjasjdadadsad")
			}
			err = tx.Commit().Error
			if err != nil {
				tx.Rollback()
				fmt.Println("dajsdjasjdjasdjsajdvj")
				return
			}

		}
	} else {
		fmt.Println("  No on owns the arena")
		return
	}

}

func Close(s socketio.Conn, req map[string]interface{}) {
	fmt.Println("Close event hit")
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
	fmt.Println("Open event hit")
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
	fmt.Println("arena reward is", arenaReward)
	res := PerkResposne{
		ArenaId:        arenaReward.ArenaId,
		PlayerId:       arenaReward.PlayerId,
		Coins:          arenaReward.Coins,
		Cash:           arenaReward.Cash,
		RepairCurrency: arenaReward.RepairCurrency,
		NextRewardTime: arenaReward.NextRewardTime.Sub(time.Now()).String(),
	}

	response.SocketResponse(utils.DATA_FETCH_SUCCESS, utils.HTTP_OK, utils.SUCCESS, res, "reward", s)

}
