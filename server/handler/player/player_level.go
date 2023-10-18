package handler

import (
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/utils"
)

func AddPlayerLevel() {
	if !utils.TableIsEmpty("player_levels") {
		multipliers := map[int]float64{
			1: 1.450,
			2: 1.250,
			3: 1.200,
			4: 1.150,
			5: 1.150,
		}

		coinsMultiplier := map[int]int{
			1: 1,
			2: 2,
			3: 3,
			4: 4,
			5: 5,
		}

		// Initialize the result slice with the first two elements
		res := []model.PlayerLevel{
			{Level: 1, XPRequired: 0, Coins: 0},
			{Level: 2, XPRequired: 100, Coins: 250},
		}

		// Loop from 2 to 49
		for i := 2; i <= 49; i++ {
			// Determine the multiplier based on the range
			rangeIndex := (i-1)/10 + 1
			multiplier := multipliers[rangeIndex]

			// Calculate the new value and round it to the nearest multiple of 5
			val := float64(res[i-1].XPRequired) * multiplier
			xp := roundToNearestMultipleOf5(val)

			// Calculate the coins based on the range
			coins := int(res[1].Coins) * coinsMultiplier[rangeIndex]

			// Append the values to the result slice
			res = append(res, model.PlayerLevel{Level: int64(i + 1), XPRequired: int64(xp), Coins: int64(coins)})
		}
		err := db.CreateRecord(&res)
		if err != nil {
			fmt.Println("Error is:", err)
			return
		}
	}

}

func roundToNearestMultipleOf5(value float64) float64 {
	return float64(int((value+2.5)/5.0) * 5)
}
