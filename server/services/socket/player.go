package socket

import (
	"fmt"
	"main/server/db"
	"main/server/response"
	"main/server/utils"

	socketio "github.com/googollee/go-socket.io"
)

var playerQuery = `
SELECT
    p.player_id,
    p.player_name,
    p.level,
    p.xp,
    p.role,
    p.email,
    p.coins,
    p.cash,
    p.repair_currency,
    (SELECT COUNT(cust_id) FROM owned_cars WHERE player_id = p.player_id) AS CarsOwned,
    COUNT(og.garage_id) AS GaragesOwned,
    (
        SELECT COUNT(arena_id) 
        FROM player_race_stats 
        WHERE player_id = p.player_id 
        AND arena_won=true
        AND win_streak > lose_streak
        AND win_streak + lose_streak = 
            CASE
                WHEN (SELECT arena_level FROM arenas WHERE arena_id = player_race_stats.arena_id) = 1 THEN 3
                WHEN (SELECT arena_level FROM arenas WHERE arena_id = player_race_stats.arena_id) = 2 THEN 5
                WHEN (SELECT arena_level FROM arenas WHERE arena_id = player_race_stats.arena_id) = 3 THEN 7
                ELSE 0
            END
    ) AS ArenaCount,
    (SELECT xp_required FROM player_levels WHERE level = p.level + 1) AS NextXPRequired,
    (SELECT xp_required FROM player_levels WHERE level = p.level) AS PrevXP,
    SUM(prh.shd_won) AS TotalShowDownWon,
    CASE
        WHEN SUM(prh.total_shd_played) > 0 THEN SUM(prh.shd_won) / SUM(prh.total_shd_played)
        ELSE 0
    END AS ShowDownWinRatio,
    SUM(prh.td_won) AS TotalTakeDownWon,
    CASE
        WHEN SUM(prh.total_td_played) > 0 THEN SUM(prh.td_won) / SUM(prh.total_td_played)
        ELSE 0
    END AS TakeDownWinRatio
FROM players p
LEFT JOIN owned_garages og ON og.player_id = p.player_id
LEFT JOIN player_race_stats prh ON prh.player_id = p.player_id
WHERE p.player_id = ?
GROUP BY
    p.player_id,
    p.player_name,
    p.level,
    p.xp,
    p.role,
    p.email,
    p.coins,
    p.cash,
    p.repair_currency;

`

func GetPlayerDetails(s socketio.Conn, req map[string]interface{}) {

	fmt.Println("Player details socket called")
	playerId := s.Context().(string)
	fmt.Println("playerid is", playerId)

	var playerResponse *response.PlayerResposne
	playerResponse, err := db.ResponseQuery(playerQuery, playerId)
	if err != nil {
		fmt.Println("error is ", err.Error())
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "ack", s)
		return
	}

	response.SocketResponse(utils.DATA_FETCH_SUCCESS, utils.HTTP_OK, utils.SUCCESS, *playerResponse, "playerDetails", s)
}

func GetPlayerDetailsCopy(playerId string) (*response.Success, error) {
	// playerId := s.Context().(string)
	fmt.Println("playerid is", playerId)

	var playerResponse *response.PlayerResposne
	playerResponse, err := db.ResponseQuery(playerQuery, playerId)
	if err != nil {
		return nil, err
	}

	resp := &response.Success{
		Status:  utils.SUCCESS,
		Code:    utils.HTTP_OK,
		Message: utils.DATA_FETCH_SUCCESS,
		Data:    *playerResponse,
	}
	return resp, nil
}
