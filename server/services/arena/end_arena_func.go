package arena

import (
	"main/server/db"
	"main/server/model"
	"main/server/request"
	"main/server/response"
	"main/server/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func UpdateArenaOwnedData(endChallReq request.EndChallengeReq, playerId2 string, ctx *gin.Context, arenaDetails model.Arena) error {
	query := "UPDATE player_race_stats SET win_time=?,arena_won=? WHERE arena_id=? AND player_id=?"
	err := db.RawExecutor(query, time.Now(), "true", endChallReq.ArenaId, playerId2)
	if err != nil {

		return err
	}

	query = "UPDATE player_race_stats SET arena_won=? WHERE arena_id=? AND player_id=?"
	err = db.RawExecutor(query, "false", endChallReq.ArenaId, endChallReq.PlayerId1)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return err
	}

	newRecord := model.ArenaReward{
		ArenaId:        endChallReq.ArenaId,
		PlayerId:       playerId2,
		Coins:          0,
		Cash:           0,
		RepairCurrency: 0,
		RewardTime:     time.Now(),
	}

	switch int64(arenaDetails.ArenaLevel) {
	case int64(utils.EASY):
		newRecord.NextRewardTime = time.Now().Add(time.Duration(utils.EASY_PERK_MINUTES) * time.Minute)
	case int64(utils.MEDIUM):
		newRecord.NextRewardTime = time.Now().Add(time.Duration(utils.MEDIUM_PERK_MINUTES) * time.Minute)
	case int64(utils.HARD):
		newRecord.NextRewardTime = time.Now().Add(time.Duration(utils.HARD_PERK_MINUTES) * time.Minute)

	}

	err = db.CreateRecord(&newRecord)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return err
	}

	query = "DELETE FROM arena_rewards WHERE player_id=? AND arena_id=?"
	err = db.RawExecutor(query, endChallReq.PlayerId1, endChallReq.ArenaId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return err
	}
	return nil
}
func BuildArenaRewardResponse(test model.RaceTypes, ctx *gin.Context, playerId2 string, endChallReq request.EndChallengeReq, status string) ([]response.RewardResponse, error) {
	var reward1, reward2 model.RaceRewards
	query := "SELECT * FROM race_rewards WHERE race_id=? AND status=?"
	err := db.QueryExecutor(query, &reward1, test.RaceId, status)
	if err != nil {
		return nil, err
	}

	_, err = EarnedRewards(playerId2, ctx, reward1)
	if err != nil {
		return nil, err
	}

	err = db.QueryExecutor(query, &reward2, endChallReq.RaceId, status)
	if err != nil {
		return nil, err
	}

	playerLevel, err := EarnedRewards(playerId2, ctx, reward2)
	if err != nil {
		return nil, err
	}
	var totalRewards = []response.RewardResponse{}

	if playerLevel != nil {

		totalRewards = append(totalRewards, response.RewardResponse{
			RewardName: "level",
			RewardData: response.RewardData{
				Coins:      playerLevel.Coins,
				Level:      playerLevel.Level,
				XPRequired: playerLevel.XPRequired,
			},
		})
	}

	totalRewards = append(totalRewards, response.RewardResponse{
		RewardName: "arena",
		RewardData: response.RewardData{
			Coins:          reward1.Coins,
			Cash:           reward1.Cash,
			RepairCurrency: reward1.RepairCurrency,
			XPGained:       reward1.XPGained,
			Status:         reward1.Status,
		},
	})

	totalRewards = append(totalRewards, response.RewardResponse{
		RewardName: "takedowns",
		RewardData: response.RewardData{
			Coins:          reward2.Coins,
			Cash:           reward2.Cash,
			RepairCurrency: reward2.RepairCurrency,
			XPGained:       reward2.XPGained,
			Status:         reward2.Status,
		},
	})
	return totalRewards, nil
}
func UpdatePlayerRaceHistory(playerId string, ctx *gin.Context, endChallReq request.EndChallengeReq, win bool) error {

	//get player race history
	var playerRaceHistory model.PlayerRaceStats
	err := db.FindById(&playerRaceHistory, playerId, utils.PLAYER_ID)
	if err != nil {
		return err
	}

	//get the details of the race type
	var raceType model.RaceTypes
	err = db.FindById(&raceType, endChallReq.RaceId, "race_id")
	if err != nil {
		return err
	}

	//update the details
	playerRaceHistory.DistanceTraveled += raceType.RaceLength
	if raceType.RaceName == "showdowns" {
		playerRaceHistory.TotalShdPlayed += 1
		if win {
			playerRaceHistory.ShdWon += 1
		}
	}
	if raceType.RaceName == "takedowns" {
		playerRaceHistory.TotalTdPlayed += 1
		if win {
			playerRaceHistory.TdWon += 1
		}
	}

	err = db.UpdateRecord(&playerRaceHistory, playerId, utils.PLAYER_ID).Error
	if err != nil {
		return err
	}
	return nil
}

func UpgradePlayerLevel(newXp int64, playerDetails *model.Player) (*model.PlayerLevel, bool, error) {
	currentLevel := playerDetails.Level
	var playerLevel model.PlayerLevel
	query := "SELECT * FROM player_levels WHERE level=?"
	err := db.QueryExecutor(query, &playerLevel, currentLevel+1)
	if err != nil {
		return nil, false, err
	}

	//give player level upgrade reward
	if newXp >= playerLevel.XPRequired {
		// Update player level
		playerDetails.Level++
		playerDetails.Coins += playerLevel.Coins

		return &playerLevel, true, nil
	}

	// fmt.Println("player Details is:", playerDetails)

	return nil, false, nil
}

func EarnedRewards(playerId string, ctx *gin.Context, rewards model.RaceRewards) (*model.PlayerLevel, error) {

	//get player details

	playerDetails, err := utils.GetPlayerDetails(playerId)
	if err != nil {
		return nil, err
	}
	//begin transaction
	tx := db.BeginTransaction()
	if tx.Error != nil {
		return nil, err
	}

	playerDetails.Coins += rewards.Coins
	playerDetails.Cash += rewards.Cash
	playerDetails.RepairCurrency += rewards.RepairCurrency
	playerDetails.XP += rewards.XPGained

	playerLevel, isUpgraded, err := UpgradePlayerLevel(playerDetails.XP, playerDetails)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	//fmt.Println("player levele is ", playerLevel)

	err = db.UpdateRecord(&playerDetails, playerId, utils.PLAYER_ID).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if isUpgraded {
		// Handle player level upgrade logic, if needed
		return playerLevel, nil
	}

	if err = tx.Commit().Error; err != nil {
		return nil, err
	}

	return nil, nil
}
