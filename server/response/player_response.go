package response

type PlayerResposne struct {
	PlayerId         string  `json:"playerId"`
	PlayerName       string  `json:"playerName" gorm:"unique"`
	Level            int     `json:"level"`
	XP               int64   `json:"xp"`
	Role             string  `json:"role"`
	Email            string  `json:"email"`
	Coins            int64   `json:"coins"`
	Cash             int64   `json:"cash"`
	RepairParts      int64   `json:"repairParts"`
	CarsOwned        int64   `json:"carsOwned"`
	GaragesOwned     int64   `json:"garagesOwned"`
	ArenasOwned      int64   `json:"arenasOwned"`
	DistanceTraveled float64 `json:"distanceTraveled"`
	XPRequired       int64   `json:"xpRequired"`
	PrevXP           int64   `json:"prevXP"`
	ShdWon           int64   `json:"showDownWon"`
	ShdWinRatio      int64   `json:"showDownWinRatio"`
	TdWon            int64   `json:"takeDownWon"`
	TdWinRatio       int64   `json:"takeDownWinRatio"`
}

type Level struct {
	LastLevelXp   int64 `json:"lastLevelXp"`
	NextLevelXp   int64 `json:"nextLevelXp"`
	CurrentReward int64 `json:"currentReward"`
	RewardValue   int64 `json:"rewardValue"`
	LevelNumber   int64 `json:"levelNumber"`
}
