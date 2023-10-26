package arena

import (
	"errors"
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/request"
	"main/server/response"
	"main/server/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func EndChallengeService(ctx *gin.Context, endChallReq request.EndChallengeReq, playerId2 string) {
	//converting the given input to format "00:00:05.1455"
	TimeInString := fmt.Sprintf("00:00:%02d.%02d%02d", int(endChallReq.Seconds), int(endChallReq.MilliSec), int(endChallReq.MicroSec))

	//fmt.Println("Time in string", TimeInString)
	winTime, err := utils.TimeConversion(TimeInString)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}
	//fmt.Println("Win time is", winTime)

	var raceType model.RaceTypes
	query := "SELECT * FROM race_types WHERE race_id=?"
	err = db.QueryExecutor(query, &raceType, endChallReq.RaceId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	if endChallReq.ArenaId == "" {
		fmt.Println("Not in arena")

		// Start a transaction
		tx := db.BeginTransaction()
		if tx.Error != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}

		defer func() {
			if err != nil {
				// Rollback the transaction if there's an error
				tx.Rollback()
				response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			}
		}()

		var oppTimeStringFormat string
		query := "SELECT time_win FROM arena_race_records WHERE player_id=? AND arena_id=? AND result='win'"
		err = db.QueryExecutor(query, &oppTimeStringFormat, endChallReq.PlayerId1, endChallReq.RaceId)
		if err != nil {
			return // Error handling will be done in the defer block
		}

		opponentTime, err := utils.TimeConversion(oppTimeStringFormat)
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
			return // Error handling will be done in the defer block
		}
		var rewards model.RaceRewards
		win := false

		query = "SELECT * FROM race_rewards WHERE race_id=? AND status=?"

		if winTime.Before(*opponentTime) {

			fmt.Println("Wins the challenge outside the arena")
			win = true

			// Check the type of the race and allot the rewards to the player
			err = db.QueryExecutor(query, rewards, endChallReq.RaceId, "win")
			if err != nil {
				return // Error handling will be done in the defer block
			}

			var totalRewards = []response.RewardResponse{}
			playerLevel, err := EarnedRewards(playerId2, ctx, rewards)
			if err != nil {
				return // Error handling will be done in the defer block
			}

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
				RewardName: raceType.RaceName,
				RewardData: response.RewardData{
					Coins:          rewards.Coins,
					Cash:           rewards.Cash,
					RepairCurrency: rewards.RepairCurrency,
					XPGained:       rewards.XPGained,
					Status:         rewards.Status,
				},
			})

			response.ShowResponse(utils.WON, utils.HTTP_OK, utils.SUCCESS, totalRewards, ctx)
			// Player wins
		} else {

			fmt.Println("player lost the challenge outside the arena")
			// Player loses
			win = false
			err = db.QueryExecutor(query, rewards, endChallReq.RaceId, "lost")
			if err != nil {
				return // Error handling will be done in the defer block
			}

			// Get player details

			var totalRewards = []response.RewardResponse{}

			totalRewards = append(totalRewards, response.RewardResponse{
				RewardName: raceType.RaceName,
				RewardData: response.RewardData{
					Coins:          rewards.Coins,
					Cash:           rewards.Cash,
					RepairCurrency: rewards.RepairCurrency,
					XPGained:       rewards.XPGained,
					Status:         rewards.Status,
				},
			})

			// Player wins
			// Give rewards to player
			playerLevel, err := EarnedRewards(playerId2, ctx, rewards)
			if err != nil {
				return // Error handling will be done in the defer block
			}

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

			response.ShowResponse(utils.WON, utils.HTTP_OK, utils.SUCCESS, totalRewards, ctx)
		}

		// Update player race history
		err = UpdatePlayerRaceHistory(playerId2, ctx, endChallReq, win)
		if err != nil {
			return // Error handling will be done in the defer block
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			return // Error handling will be done in the defer block
		}
	} else {
		fmt.Println("Taking challenge in arena")
		//player is taking challenge in arena

		var winCount int64
		query := "SELECT win_streak FROM player_race_stats WHERE arena_id=? AND player_id=?"
		err := db.QueryExecutor(query, &winCount, endChallReq.ArenaId, playerId2)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			winCount = 0
		}
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
			return
		}

		//repeated code correct it
		var lostCount int64
		query = "SELECT lose_streak FROM player_race_stats WHERE arena_id=? AND player_id=?"
		err = db.QueryExecutor(query, &lostCount, endChallReq.ArenaId, playerId2)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			lostCount = 0
		}
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
			return
		}

		var oppTimeStringFormat []string
		query = "SELECT time_win FROM arena_race_records WHERE player_id=? AND arena_id=? ORDER BY created_at"
		err = db.QueryExecutor(query, &oppTimeStringFormat, endChallReq.PlayerId1, endChallReq.ArenaId)
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}

		fmt.Println("Opponent string format is ", oppTimeStringFormat)
		if oppTimeStringFormat == nil {
			response.ShowResponse("There is nothing in that database for That playerID sent in body ", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
			return
		}
		var arenaDetails model.Arena
		err = db.FindById(&arenaDetails, endChallReq.ArenaId, "arena_id")
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}

		var maxRaces int64

		//check for the win count
		switch arenaDetails.ArenaLevel {
		case int64(utils.EASY):
			maxRaces = utils.EASY_ARENA_SLOT
		case int64(utils.MEDIUM):
			maxRaces = utils.MEDIUM_ARENA_SLOT
		case int64(utils.HARD):
			maxRaces = utils.HARD_ARENA_SLOT

		}

		fmt.Println("Total race count ", winCount+lostCount)
		fmt.Println("Max race cont", maxRaces)
		if (winCount + lostCount) == maxRaces {

			if winCount > lostCount {
				response.ShowResponse("Player already Won the arena", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
				return
			} else if lostCount > winCount {

				//reset the arena series and empty the temp record for that player
				query = "UPDATE player_race_stats SET win_streak=0 AND lose_streak=0 WHERE arena_id=? AND player_id=?"
				err = db.RawExecutor(query, endChallReq.ArenaId, playerId2)
				if err != nil {
					response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
					return
				}
				query = "DELETE FROM temp_race_records WHERE arena_id=? AND player_id=?"
				err = db.RawExecutor(query, endChallReq.ArenaId, playerId2)
				if err != nil {
					response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
					return
				}

			}

		}
		opponentTime, err := utils.TimeConversion(oppTimeStringFormat[(winCount + lostCount)])
		if err != nil {

			response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
			return
		}
		fmt.Println("Compared opponent time is:", *opponentTime)
		if winTime.Before(*opponentTime) {

			fmt.Println("Compared opponent time is:", opponentTime)
			fmt.Println("player won in arena")
			//player wins the a series in arena
			//add the count to arenaRaceWins
			var exists bool
			query := "SELECT EXISTS (SELECT * FROM player_race_stats WHERE arena_id=? AND player_id=?)"
			err := db.QueryExecutor(query, &exists, endChallReq.ArenaId, playerId2)
			if err != nil {
				response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
				return
			}
			if exists {
				//increment the number of wins
				query := "UPDATE player_race_stats SET win_streak=win_streak+1 WHERE  arena_id=? AND player_id=?"
				err := db.RawExecutor(query, endChallReq.ArenaId, playerId2)
				if err != nil {
					response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
					return
				}

				//check if the player is eligible for arena reward or not
				var arenaSeries model.PlayerRaceStats
				query = "SELECT * FROM player_race_stats WHERE arena_id=? AND player_id=?"
				err = db.QueryExecutor(query, &arenaSeries, endChallReq.ArenaId, playerId2)
				if err != nil {
					response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
					return
				}

				//change variable name
				var test model.RaceTypes
				query = "SELECT * FROM race_types WHERE race_level=? AND race_name='arena'"
				err = db.QueryExecutor(query, &test, arenaDetails.ArenaLevel)
				if err != nil {
					response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
					return
				}
				fmt.Println("test is:", test)
				fmt.Println("arena series", arenaSeries)
				if (arenaSeries.WinStreak+arenaSeries.LoseStreak) == test.RaceSeries && (arenaSeries.WinStreak > arenaSeries.LoseStreak) {
					//player won the arena
					fmt.Println("wining the arena")

					//as reward1 donot contain xp so playerLevel field is omoited here
					totalRewards, err := BuildArenaRewardResponse(test, ctx, playerId2, endChallReq, "win")
					if err != nil {

						response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
						return
					}

					tempRecord := &model.TempRaceRecords{
						PlayerId: playerId2,
						ArenaId:  endChallReq.ArenaId,
						TimeWin:  TimeInString,
						CustId:   endChallReq.CustId,
						Result:   "win",
					}
					err = db.CreateRecord(&tempRecord)
					if err != nil {
						response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
						return
					}

					//get from temp table and update it into original table
					var tempRecords []model.TempRaceRecords
					query = "SELECT * FROM temp_race_records WHERE player_id=? AND arena_id=?"
					err = db.QueryExecutor(query, &tempRecords, playerId2, endChallReq.ArenaId)
					if err != nil {
						response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
						return
					}

					fmt.Println("TempRecords are", tempRecords)

					//update is not working so first deleteing those reords and then add new records

					query = "DELETE FROM  arena_race_records  WHERE player_id=? AND arena_id=?"
					err = db.RawExecutor(query, endChallReq.PlayerId1, endChallReq.ArenaId)
					if err != nil {
						response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
						return
					}

					for _, rec := range tempRecords {

						newRecord := model.ArenaRaceRecord{
							PlayerId: playerId2,
							ArenaId:  endChallReq.ArenaId,
							TimeWin:  rec.TimeWin,
							Result:   rec.Result,
							CustId:   rec.CustId,
						}

						err = db.CreateRecord(&newRecord)
						if err != nil {
							response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
							return
						}
					}

					query = "DELETE FROM  temp_race_records  WHERE player_id=? AND arena_id=?"
					err = db.RawExecutor(query, playerId2, endChallReq.ArenaId)
					if err != nil {
						response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
						return
					}

					err = UpdatePlayerRaceHistory(playerId2, ctx, endChallReq, true)
					if err != nil {
						response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
						return
					}

					//give both rewards arena and takedown
					response.ShowResponse(utils.WON, utils.HTTP_OK, utils.SUCCESS, totalRewards, ctx)

					//add a 24 hour timer after the arena is won
					///if after the 24 hour there is no entery in carSlots table then the arebna will be given back to the AI

					//it should be 24 hours

				} else if (arenaSeries.WinStreak+arenaSeries.LoseStreak) == test.RaceSeries && (arenaSeries.WinStreak < arenaSeries.LoseStreak) {
					//player lost the arena so give arena lost reward
					fmt.Println("Lost the arena")

					totalRewards, err := BuildArenaRewardResponse(test, ctx, playerId2, endChallReq, "lost")
					if err != nil {

						response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
						return
					}
					err = UpdatePlayerRaceHistory(playerId2, ctx, endChallReq, true)
					if err != nil {
						response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
						return
					}
					response.ShowResponse(utils.LOSE, utils.HTTP_OK, utils.SUCCESS, totalRewards, ctx)
				} else {

					//player won the challenge but not the arena
					var reward model.RaceRewards
					var totalRewards = []response.RewardResponse{}
					query = "SELECT * FROM race_rewards WHERE race_id=? AND status=?"
					err = db.QueryExecutor(query, &reward, endChallReq.RaceId, "win")
					if err != nil {
						response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
						return
					}
					playerLevel, err := EarnedRewards(playerId2, ctx, reward)
					if err != nil {
						response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
						return
					}

					//create a temp database for storing the result of the win
					tempRecord := &model.TempRaceRecords{
						PlayerId: playerId2,
						ArenaId:  endChallReq.ArenaId,
						TimeWin:  TimeInString,
						Result:   "win",
						CustId:   endChallReq.CustId,
					}
					err = db.CreateRecord(&tempRecord)
					if err != nil {
						response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
						return
					}

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
						RewardName: raceType.RaceName,
						RewardData: response.RewardData{
							Coins:          reward.Coins,
							Cash:           reward.Cash,
							RepairCurrency: reward.RepairCurrency,
							XPGained:       reward.XPGained,
							Status:         reward.Status,
						},
					})

					err = UpdatePlayerRaceHistory(playerId2, ctx, endChallReq, true)
					if err != nil {
						response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
						return
					}
					response.ShowResponse(utils.WON, utils.HTTP_OK, utils.SUCCESS, totalRewards, ctx)
				}
			} else {

				//test it this else might create error

				//create a record and set the initail win to 1
				// arenaSeriesRecord := model.ArenaSeries{
				// 	ArenaId:    endChallReq.ArenaId,
				// 	PlayerId:   playerId2,
				// 	WinStreak:  1,
				// 	LoseStreak: 0,
				// }

				arenaSeriesRecord := model.PlayerRaceStats{
					PlayerId:   playerId2,
					ArenaId:    &endChallReq.ArenaId,
					WinStreak:  1,
					LoseStreak: 0,
				}

				err := db.CreateRecord(&arenaSeriesRecord)
				if err != nil {
					response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
					return
				}

				var reward model.RaceRewards
				var totalRewards = []response.RewardResponse{}
				query = "SELECT * FROM race_rewards WHERE race_id=? AND status=?"
				err = db.QueryExecutor(query, &reward, endChallReq.RaceId, "win")
				if err != nil {
					response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
					return
				}

				playerLevel, err := EarnedRewards(playerId2, ctx, reward)
				if err != nil {
					response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
					return
				}

				tempRecord := &model.TempRaceRecords{
					PlayerId: playerId2,
					ArenaId:  endChallReq.ArenaId,
					TimeWin:  TimeInString,
					Result:   "win",
					CustId:   endChallReq.CustId,
				}
				err = db.CreateRecord(&tempRecord)
				if err != nil {
					response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
					return
				}
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
					RewardName: raceType.RaceName,
					RewardData: response.RewardData{
						Coins:          reward.Coins,
						Cash:           reward.Cash,
						RepairCurrency: reward.RepairCurrency,
						XPGained:       reward.XPGained,
						Status:         reward.Status,
					},
				})
				err = UpdatePlayerRaceHistory(playerId2, ctx, endChallReq, true)
				if err != nil {
					response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
					return
				}
				response.ShowResponse(utils.WON, utils.HTTP_OK, utils.SUCCESS, totalRewards, ctx)

			}

		} else {
			// fghddfhdfghdfhdfghdfgh
			//check that if the record exists in arena series or not
			fmt.Println("Player lost the race...")
			var exists bool
			query := "SELECT EXISTS (SELECT * FROM player_race_stats WHERE arena_id=? AND player_id=?)"
			err := db.QueryExecutor(query, &exists, endChallReq.ArenaId, playerId2)
			if err != nil {
				response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
				return
			}
			if exists {

				query := "UPDATE player_race_stats SET lose_streak=lose_streak+1 WHERE  arena_id=? AND player_id=?"
				err := db.RawExecutor(query, endChallReq.ArenaId, playerId2)
				if err != nil {
					response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
					return
				}
			}
			// } else {
			// 	arenaSeriesRecord := model.ArenaSeries{
			// 		ArenaId:    endChallReq.ArenaId,
			// 		PlayerId:   playerId2,
			// 		WinStreak:  0,
			// 		LoseStreak: 1,
			// 	}
			// 	err := db.CreateRecord(&arenaSeriesRecord)
			// 	if err != nil {
			// 		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			// 		return
			// 	}
			// }

			//player Lost
			var reward model.RaceRewards
			var totalRewards = []response.RewardResponse{}
			query = "SELECT * FROM race_rewards WHERE race_id=? AND status=?"
			err = db.QueryExecutor(query, &reward, endChallReq.RaceId, "lost")
			if err != nil {
				response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
				return
			}
			tempRecord := &model.TempRaceRecords{
				PlayerId: playerId2,
				ArenaId:  endChallReq.ArenaId,
				TimeWin:  TimeInString,
				Result:   "lost",
				CustId:   endChallReq.CustId,
			}
			err = db.CreateRecord(&tempRecord)
			if err != nil {
				response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
				return
			}
			playerLevel, err := EarnedRewards(playerId2, ctx, reward)
			if err != nil {
				response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
				return
			}

			//the race is the last race of that arena then give arena lost reward also

			//check if the player is eligible for arena reward or not
			var arenaSeries model.PlayerRaceStats
			query = "SELECT * FROM player_race_stats WHERE arena_id=? AND player_id=?"
			err = db.QueryExecutor(query, &arenaSeries, endChallReq.ArenaId, playerId2)
			if err != nil {
				response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
				return
			}

			//change variable name
			var test model.RaceTypes
			query = "SELECT * FROM race_types WHERE race_level=? AND race_name='arena'"
			err = db.QueryExecutor(query, &test, arenaDetails.ArenaLevel)
			if err != nil {
				response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
				return
			}

			if arenaSeries.LoseStreak+arenaSeries.WinStreak == test.RaceSeries && arenaSeries.LoseStreak > arenaSeries.WinStreak {

				fmt.Println("Player lost the arena and also lost the last race")
				var reward1 model.RaceRewards
				query := "SELECT * FROM race_rewards WHERE race_id=? AND status=?"
				err := db.QueryExecutor(query, &reward1, test.RaceId, "lost")
				if err != nil {
					response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)

					return
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

			} else if arenaSeries.LoseStreak+arenaSeries.WinStreak == test.RaceSeries && arenaSeries.LoseStreak < arenaSeries.WinStreak {
				fmt.Println("Player won the arena but lost the last race")

				query = "UPDATE player_race_stats SET player_id=? WHERE player_id=? AND arena_id=?"
				err = db.RawExecutor(query, playerId2, endChallReq.PlayerId1, endChallReq.ArenaId)
				if err != nil {
					response.ShowResponse("Unable to update the owned arena details", utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
					return
				}

				//get from temp table and update it into original table
				var tempRecords []model.TempRaceRecords
				query = "SELECT * FROM temp_race_records WHERE player_id=? AND arena_id=?"
				err = db.QueryExecutor(query, &tempRecords, playerId2, endChallReq.ArenaId)
				if err != nil {
					response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
					return
				}

				fmt.Println("TempRecords are", tempRecords)

				//update is not working so first deleteing those reords and then add new records

				query = "DELETE FROM  arena_race_records  WHERE player_id=? AND arena_id=?"
				err = db.RawExecutor(query, endChallReq.PlayerId1, endChallReq.ArenaId)
				if err != nil {
					response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
					return
				}

				for _, rec := range tempRecords {

					newRecord := model.ArenaRaceRecord{
						PlayerId: playerId2,
						ArenaId:  endChallReq.ArenaId,
						TimeWin:  rec.TimeWin,
						Result:   rec.Result,
						CustId:   rec.CustId,
					}

					err = db.CreateRecord(&newRecord)
					if err != nil {
						response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
						return
					}
				}
				query = "DELETE FROM  temp_race_records  WHERE player_id=? AND arena_id=?"
				err = db.RawExecutor(query, playerId2, endChallReq.ArenaId)
				if err != nil {
					response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
					return
				}

				var reward1 model.RaceRewards
				query := "SELECT * FROM race_rewards WHERE race_id=? AND status=?"

				err := db.QueryExecutor(query, &reward1, test.RaceId, "win")
				if err != nil {
					response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)

					return
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
			}

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
				RewardName: raceType.RaceName,
				RewardData: response.RewardData{
					Coins:          reward.Coins,
					Cash:           reward.Cash,
					RepairCurrency: reward.RepairCurrency,
					XPGained:       reward.XPGained,
					Status:         reward.Status,
				},
			})
			err = UpdatePlayerRaceHistory(playerId2, ctx, endChallReq, false)
			if err != nil {
				response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
				return
			}
			response.ShowResponse(utils.LOSE, utils.HTTP_OK, utils.SUCCESS, totalRewards, ctx)
		}

	}

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
func GetArenaOwnerService(ctx *gin.Context, arenaId string) {
	// Check if the arena exists in the database
	if !db.RecordExist("arenas", arenaId, "arena_id") {
		response.ShowResponse(utils.ARENA_NOT_FOUND, utils.HTTP_NOT_FOUND, utils.FAILURE, nil, ctx)
		return
	}

	// Declare the owner model to store player race stats
	var owner model.PlayerRaceStats

	// Query to get the player with the highest win streak
	query := "SELECT * FROM player_race_stats WHERE arena_id=? and win_streak>lose_streak ORDER BY updated_at DESC LIMIT 1"

	// Execute the query and store the result in 'owner'
	err := db.QueryExecutor(query, &owner, arenaId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	// Declare a model for player details
	var playerDetails model.Player

	// Query to get player details using 'owner.PlayerId'
	query = "SELECT * FROM players WHERE player_id=?"

	// Execute the query and store the result in 'playerDetails'
	err = db.QueryExecutor(query, &playerDetails, owner.PlayerId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	// Declare a slice to store arena race records
	var arenaRecord []string

	// Query to get the 'time_win' from arena race records
	query = "SELECT time_win from arena_race_records WHERE arena_id=? AND player_id=? ORDER BY created_at"

	// Execute the query and store the result in 'arenaRecord'
	err = db.QueryExecutor(query, &arenaRecord, arenaId, playerDetails.PlayerId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	// Declare a slice to store car information

	var carStruct2 []response.CarCustom

	// Declare a response struct to format the final response
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
		Cars []response.CarRes `json:"cars"`
	}

	// Populate the response struct with player and arena details
	resp.PlayerId = playerDetails.PlayerId
	resp.PlayerName = playerDetails.PlayerName
	resp.ArenaId = arenaId

	// Parse and format the arena record times
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

	// Query to fetch car customizations associated with the player's records

	query = `SELECT pcc.*,c.class,c.car_name
        FROM player_car_customisations pcc
        JOIN cars c ON c.car_id=pcc.car_id 
        JOIN arena_race_records arr ON arr.cust_id=pcc.cust_id
        WHERE arr.arena_id=? AND arr.player_id=?;`

	// Execute the query and store the result in 'carStruct2'
	err = db.QueryExecutor(query, &carStruct2, arenaId, playerDetails.PlayerId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	// Append carStruct2 to the 'resp.Cars' slice

	for _, details := range carStruct2 {
		carRes := response.CarRes{
			CustId:  details.CustId,
			CarId:   details.CarId,
			CarName: details.CarName,
			Rarity:  details.Class,
		}

		carCustomise, _ := utils.CustomiseMapping(details.CustId, "player_car_customisations")
		carRes.CarCurrentData.Customization = *carCustomise
		carRes.CarCurrentData.Stats.Power = details.Power
		carRes.CarCurrentData.Stats.Grip = details.Grip
		carRes.CarCurrentData.Stats.Weight = details.Weight
		carRes.CarCurrentData.Stats.ShiftTime = details.ShiftTime
		carRes.CarCurrentData.Stats.OVR = details.OVR
		carRes.CarCurrentData.Stats.Durability = details.Durability
		carRes.Status.Owned = true
		resp.Cars = append(resp.Cars, carRes)

	}

	// Show the final response with success message
	response.ShowResponse(utils.DATA_FETCH_SUCCESS, utils.HTTP_OK, utils.SUCCESS, resp, ctx)
}

func EnterArenaService(ctx *gin.Context, enterReq request.GetArenaReq, playerId string) {
	// Check if the arena is already owned by the player
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM player_race_stats WHERE arena_id = ? AND player_id = ? AND win_streak > lose_streak)"
	err := db.QueryExecutor(query, &exists, enterReq.ArenaId, playerId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	// If the player already owns the arena, show an error response
	if exists {
		response.ShowResponse(utils.ARENA_ALREADY_OWNED, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	// Check if the player has any previous records in the arena
	query = "SELECT EXISTS (SELECT 1 FROM player_race_stats WHERE arena_id = ? AND player_id = ?)"
	err = db.QueryExecutor(query, &exists, enterReq.ArenaId, playerId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	// If the player has previous records in the arena, reset their statistics
	if exists {
		// Reset the win and lose streak to 0
		query := "UPDATE player_race_stats SET win_streak = 0, lose_streak = 0 WHERE arena_id = ? AND player_id = ?"
		err = db.RawExecutor(query, enterReq.ArenaId, playerId)
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}

		// Delete temporary race records for the player in the arena
		query = "DELETE FROM temp_race_records WHERE player_id = ? AND arena_id = ?"
		err = db.RawExecutor(query, playerId, enterReq.ArenaId)
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
			return
		}
	}

	// Show a success response if all checks and updates are completed
	response.ShowResponse("Successfully", utils.HTTP_OK, utils.SUCCESS, nil, ctx)
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
	if utils.IsExisting(query, playerId, replaceReq.ExistingCustId) {
		response.ShowResponse(utils.CAR_ALREADY_ALLOTTED, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	// Replace the car in the slot
	query = "UPDATE arena_cars SET cust_id = ? WHERE player_id = ? AND arena_id = ? AND cust_id="
	err := db.RawExecutor(query, replaceReq.NewCustId, playerId, replaceReq.ArenaId, replaceReq.ExistingCustId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	response.ShowResponse(utils.CAR_REPLACED_SUCCESS, utils.HTTP_OK, utils.SUCCESS, replaceReq, ctx)
}

func ArenaCarService(ctx *gin.Context, playerId string) {
	var temp1 []response.CarCustom
	query := `SELECT pc.*,c.class,c.car_name
				FROM player_car_customisations pc
				JOIN owned_cars oc ON pc.cust_id = oc.cust_id
				LEFT JOIN arena_cars ac ON pc.cust_id = ac.cust_id
				LEFT JOIN cars c ON pc.car_id=c.car_id
				WHERE oc.player_id = ?
				AND ac.cust_id IS NULL;`

	err := db.QueryExecutor(query, &temp1, playerId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	var res []response.CarRes
	for _, t := range temp1 {

		record := &response.CarRes{
			CustId:  t.CustId,
			CarId:   t.CarId,
			CarName: t.CarName,
			Rarity:  t.Class,
		}
		record.Status.Owned = true
		record.Status.Purchasable = false
		record.CarCurrentData.Stats.Power = t.Power
		record.CarCurrentData.Stats.Grip = t.Grip
		record.CarCurrentData.Stats.ShiftTime = t.ShiftTime
		record.CarCurrentData.Stats.Weight = t.Weight
		record.CarCurrentData.Stats.OVR = t.OVR
		record.CarCurrentData.Stats.Durability = t.Durability
		record.CarCurrentData.Stats.NitrousTime = t.NitrousTime

		carCustomise, _ := utils.CustomiseMapping(t.CustId, "player_car_customisations")

		record.CarCurrentData.Customization = *carCustomise
		res = append(res, *record)
	}

	response.ShowResponse(utils.DATA_FETCH_SUCCESS, utils.HTTP_OK, utils.SUCCESS, res, ctx)

}

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

	var arenaWinTime time.Time
	query = "SELECT win_time FROM owned_battle_arenas WHERE arena_id=?"
	err = db.QueryExecutor(query, &arenaWinTime, arenaId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	temp := time.Since(arenaWinTime.Add(24 * time.Hour))
	arenaSlotData.ArenaWinTime = temp.Abs().String()

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
