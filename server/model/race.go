package model

import "gorm.io/gorm"

type RaceTypes struct {
	RaceId      int64       `json:"raceId" gorm:"autoIncrement;primaryKey"`
	RaceRewards RaceRewards `json:"raceRewards" gorm:"references:RaceId;foreignKey:RaceId;constraint:OnDelete:CASCADE"`
	RaceName    string      `json:"raceName"`
	RaceLength  float64     `json:"raceLength"`
	RaceSeries  int64       `json:"raceSeries"`
	RaceLevel   int64       `json:"raceLevel"`
}

type RaceRewards struct {
	gorm.Model
	RaceId int64 `json:"raceId"`
	//RaceTypes   RaceTypes `json:"-" gorm:"references:RaceId;constraint:OnDelete:CASCADE"`
	Coins       int64  `json:"coins"`
	Cash        int64  `json:"cash"`
	RepairParts int64  `json:"repairParts"`
	Status      string `json:"status"`
	XPGained    int64  `json:"xpGained"`
}
