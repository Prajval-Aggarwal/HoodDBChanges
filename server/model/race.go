package model

import (
	"time"

	"gorm.io/gorm"
)

type RaceTypes struct {
	RaceId      int64       `json:"raceId" gorm:"autoIncrement;primaryKey"`
	RaceRewards RaceRewards `json:"raceRewards" gorm:"references:RaceId;foreignKey:RaceId;constraint:OnDelete:CASCADE"`
	RaceName    string      `json:"raceName"`
	RaceLength  float64     `json:"raceLength"`
	RaceSeries  int64       `json:"raceSeries"`
	RaceLevel   int64       `json:"raceLevel"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

type RaceRewards struct {
	Id             string `json:"id" gorm:"unique;default:uuid_generate_v4();primaryKey,omitempty"`
	RaceId         int64  `json:"raceId"`
	Coins          int64  `json:"coins"`
	Cash           int64  `json:"cash"`
	RepairCurrency int64  `json:"repairParts"`
	Status         string `json:"status"`
	XPGained       int64  `json:"xpGained"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt
}
