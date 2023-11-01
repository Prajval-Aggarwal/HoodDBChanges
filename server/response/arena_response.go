package response

import (
	"main/server/model"
	"time"
)

type ArenaResponse struct {
	ArenaId       string                `json:"arenaId"`
	ArenaName     string                `json:"arenaName"`
	ArenaLevel    int64                 `json:"arenaLevel"`
	Longitude     float64               `json:"longitude"`
	Latitude      float64               `json:"latitude"`
	CreatedAt     time.Time             `json:"createdAt,omitempty"`
	RewardData    model.ArenaLevelPerks `json:"rewardData"`
	NumberOfRaces int64                 `json:"numberOfRaces"`
	RewardTime    string                `json:"rewardTimer"`
}

type ArenaSlotResponse struct {
	RewardData    model.ArenaLevelPerks `json:"rewardData"`
	ArenaWinTime  string                `json:"slotAddTime"`
	ArenaPerkTime string                `json:"rewardTime"`
	TotalSlots    int                   `json:"totalSlots"`
	CarDetails    []CarRes              `json:"carDetails"`
}
