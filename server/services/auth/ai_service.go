package auth

import (
	"log"
	"main/server/db"
	"main/server/model"

	"github.com/google/uuid"
)

func AddAiToDB() {
	slice := []string{
		"sophia",
		"ethan",
		"olivia",
		"liam",
		"ava",
		"noah",
		"mia",
		"jackson",
		"isabella",
		"aiden",
	}

	//create record in player table as user type of ai

	// Generate a new player UUID and access token expiration time (48 hours from now).

	// Create a new player record with default values.

	var count int64
	query := "SELECT COUNT(*) FROM players"
	err := db.QueryExecutor(query, &count)
	if err != nil {
		return
	}

	if count == 0 {

		for i := 0; i < len(slice); i++ {
			playerUUID := uuid.New().String()
			playerRecord := model.Player{
				PlayerId:       playerUUID,
				PlayerName:     slice[i],
				Level:          1,
				Role:           "ai",
				Coins:          100000,
				Cash:           100000,
				RepairCurrency: 0,
			}

			// Create a new access token with player claims and expiration time.

			// Create player and race history records in the database.
			err := db.CreateRecord(&playerRecord)
			if err != nil {
				log.Fatal(err.Error())
				return
			}
			playerRaceHist := model.PlayerRaceStats{
				PlayerId:         playerUUID,
				ArenaId:          nil,
				WinStreak:        0,
				LoseStreak:       0,
				DistanceTraveled: 0,
				ShdWon:           0,
				TotalShdPlayed:   0,
				TdWon:            0,
				TotalTdPlayed:    0,
			}
			err = db.CreateRecord(&playerRaceHist)
			if err != nil {
				log.Fatal(err.Error())

				return
			}
		}
	}

}
