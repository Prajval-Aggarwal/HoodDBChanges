package model

type RaceTypes struct {
	RaceId     string  `json:"raceId" gorm:"autoIncrement;primaryKey"`
	RaceName   string  `json:"raceName"`
	RaceLength float64 `json:"raceLength"`
	RaceSeries int64   `json:"raceSeries"`
	RaceLevel  int64   `json:"raceLevel"`
}

type RaceRewards struct {
	Id          string    `json:"id" gorm:"unique;default:uuid_generate_v4();primaryKey,omitempty"`
	RaceId      string    `json:"race_id"`
	RaceTypes   RaceTypes `json:"-" gorm:"references:RaceId;constraint:OnDelete:CASCADE"`
	Coins       int64     `json:"coins"`
	Cash        int64     `json:"cash"`
	RepairParts int64     `json:"repairParts"`
	Status      string    `json:"status"`
	XPGained    int64     `json:"xpGained"`
}
