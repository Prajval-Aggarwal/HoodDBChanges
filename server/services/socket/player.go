package socket

import (
	"fmt"
	"main/server/db"
	"main/server/response"
	"main/server/utils"

	socketio "github.com/googollee/go-socket.io"
)

func GetPlayerDetails(s socketio.Conn, req map[string]interface{}) {
	playerId := s.Context().(string)
	fmt.Println("playerid is", playerId)

	var playerResponse *response.PlayerResposne
	query := `SELECT
    p.player_id,
    p.player_name,
    p.level,
    p.xp,
    p.role,
    p.email,
    p.coins,
    p.cash,
    p.repair_currency,
    COUNT(oc.cust_id) AS CarsOwned,
    COUNT(og.garage_id) AS GaragesOwned,
 	(SELECT COUNT(arena_id) FROM player_race_stats WHERE player_id = ? AND win_streak > lose_streak) AS ArenaCount,
    prh.distance_traveled,
    pl.xp_required AS NextXPRequired,    
    CASE
        WHEN p.level = 1 THEN 0 
        ELSE plPrev.xp_required
    END AS PrevXP,
    prh.shd_won AS ShowDownWon,
    CASE
        WHEN prh.total_shd_played > 0 THEN prh.shd_won / prh.total_shd_played    ELSE 0
    END AS ShowDownWinRatio,
    prh.td_won AS TakeDownWon,
    CASE
        WHEN prh.total_td_played > 0 THEN prh.td_won / prh.total_td_played
        ELSE 0
    END AS TakeDownWinRatio
	FROM players p
	LEFT JOIN owned_cars oc ON oc.player_id = p.player_id
	LEFT JOIN owned_garages og ON og.player_id = p.player_id
	LEFT JOIN player_race_stats prh ON prh.player_id = p.player_id
	LEFT JOIN player_levels pl ON pl.level = p.level + 1
	LEFT JOIN player_levels plPrev ON plPrev.level = p.level
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
		p.repair_currency,
		prh.distance_traveled,
		prh.shd_won,
		prh.total_shd_played,
		prh.td_won,
		prh.total_td_played,
		pl.xp_required,
		plPrev.xp_required;`
	playerResponse, err := db.ResponseQuery(query, playerId, playerId)
	if err != nil {
		fmt.Println("error is ", err.Error())
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "ack", s)
		return
	}

	response.SocketResponse(utils.DATA_FETCH_SUCCESS, utils.HTTP_OK, utils.SUCCESS, *playerResponse, "playerDetails", s)
}

func GetPlayerDetailsCopy(playerId string) (*response.PlayerResposne, error) {
	// playerId := s.Context().(string)
	fmt.Println("playerid is", playerId)

	var playerResponse *response.PlayerResposne
	query := `SELECT
    p.player_id,
    p.player_name,
    p.level,
    p.xp,
    p.role,
    p.email,
    p.coins,
    p.cash,
    p.repair_currency,
    COUNT(oc.cust_id) AS CarsOwned,
    COUNT(og.garage_id) AS GaragesOwned,
 	(SELECT COUNT(arena_id) FROM player_race_stats WHERE player_id = ? AND win_streak > lose_streak) AS ArenaCount,
    prh.distance_traveled,
    pl.xp_required AS NextXPRequired,    
    CASE
        WHEN p.level = 1 THEN 0 
        ELSE plPrev.xp_required
    END AS PrevXP,
    prh.shd_won AS ShowDownWon,
    CASE
        WHEN prh.total_shd_played > 0 THEN prh.shd_won / prh.total_shd_played    ELSE 0
    END AS ShowDownWinRatio,
    prh.td_won AS TakeDownWon,
    CASE
        WHEN prh.total_td_played > 0 THEN prh.td_won / prh.total_td_played
        ELSE 0
    END AS TakeDownWinRatio
	FROM players p
	LEFT JOIN owned_cars oc ON oc.player_id = p.player_id
	LEFT JOIN owned_garages og ON og.player_id = p.player_id
	LEFT JOIN player_race_stats prh ON prh.player_id = p.player_id
	LEFT JOIN player_levels pl ON pl.level = p.level + 1
	LEFT JOIN player_levels plPrev ON plPrev.level = p.level
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
		p.repair_currency,
		prh.distance_traveled,
		prh.shd_won,
		prh.total_shd_played,
		prh.td_won,
		prh.total_td_played,
		pl.xp_required,
		plPrev.xp_required;`
	playerResponse, err := db.ResponseQuery(query, playerId, playerId)
	if err != nil {
		return nil, err
	}

	return playerResponse, nil
}
