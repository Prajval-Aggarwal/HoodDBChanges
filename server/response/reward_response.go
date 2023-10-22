package response

type RewardData struct {
	Coins          int64  `json:"coins,omitempty"`
	Cash           int64  `json:"cash,omitempty"`
	RepairCurrency int64  `json:"repairParts,omitempty"`
	XPGained       int64  `json:"xpGained,omitempty"`
	Status         string `json:"status,omitempty"`
	Level          int64  `json:"level,omitempty"`
	XPRequired     int64  `json:"xpRequired,omitempty"`
}

type RewardResponse struct {
	RewardName string     `json:"rewardName"`
	RewardData RewardData `json:"rewardData"`
}
