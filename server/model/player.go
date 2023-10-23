package model

import (
	"time"

	"gorm.io/gorm"
)

type Player struct {
	PlayerId       string `json:"playerId,omitempty" gorm:"unique;primaryKey"`
	PlayerName     string `json:"playerName" gorm:"unique"`
	Level          int64  `json:"level,omitempty"`
	Role           string `json:"role,omitempty"`
	XP             int64  `json:"xp,omitempty"`
	Email          string `json:"email,omitempty"`
	Coins          int64  `json:"coins,omitempty"`
	Cash           int64  `json:"cash,omitempty"`
	RepairCurrency int64  `json:"repairParts,omitempty"`
	// DeviceId    string `json:"deviceId,omitempty"`
	OS        int64 `json:"os,omitempty"` // o for android 1 for ios
	CreatedAt time.Time
	UpdatedAt time.Time `gorm:"autoUpdateTime:true"`
	DeletedAt gorm.DeletedAt
}

type OwnedCars struct {
	Id        string `json:"id" gorm:"unique;default:uuid_generate_v4();primaryKey,omitempty"`
	PlayerId  string `json:"player_Id"`
	Player    Player `json:"-" gorm:"references:PlayerId;constraint:OnDelete:CASCADE"`
	Selected  bool   `json:"selected"`
	CustId    string `json:"custId" `
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
	// PlayerCarCustomisation PlayerCarCustomisation `json:"-" gorm:"references:CustId;foreignKey:CustId;constraint:OnDelete:CASCADE"`
}

type OwnedGarage struct {
	Id        string `json:"id" gorm:"unique;default:uuid_generate_v4();primaryKey,omitempty"`
	PlayerId  string `json:"playerId,omitempty"`
	Player    Player `json:"-" gorm:"references:PlayerId;constraint:OnDelete:CASCADE"`
	GarageId  string `json:"garageId,omitempty"`
	Garage    Garage `json:"-" gorm:"references:GarageId;constraint:OnDelete:CASCADE"`
	CarLimit  int64  `json:"carLimit,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type PlayerCarCustomisation struct {
	CustId            string          `json:"custId"  gorm:"unique;default:uuid_generate_v4();primaryKey,omitempty"`
	OwnedCars         OwnedCars       `json:"ownedCars" gorm:"references:CustId;foreignKey:CustId;constraint:OnDelete:CASCADE"` // Making cust id forign key in owned cars table
	GarageCars        GarageCars      `json:"garageCars" gorm:"references:CustId;foreignKey:CustId;constraint:OnDelete:CASCADE"`
	ArenaCars         ArenaCars       `json:"arenaCars" gorm:"references:CustId;foreignKey:CustId;constraint:OnDelete:CASCADE"`
	ArenaRaceRecord   ArenaRaceRecord `json:"arenaRaceRecords" gorm:"references:CustId;foreignKey:CustId;constraint:OnDelete:CASCADE"`
	PlayerId          string          `json:"playerId,omitempty"`
	Player            Player          `json:"-" gorm:"references:PlayerId;constraint:OnDelete:CASCADE"`
	CarId             string          `json:"carId,omitempty"`
	Car               Car             `json:"-" gorm:"references:CarId;constraint:OnDelete:CASCADE"`
	CarLevel          int64           `json:"carLevel,omitempty"`
	Power             int64           `json:"power,omitempty"`
	Grip              int64           `json:"grip,omitempty"`
	ShiftTime         float64         `json:"shiftTime,omitempty"`
	Weight            int64           `json:"weight,omitempty"`
	OVR               float64         `json:"or,omitempty"` //overall rating of the car
	Durability        int64           `json:"Durability,omitempty"`
	NitrousTime       float64         `json:"nitrousTime,omitempty"` //increased when nitrous is upgraded
	ColorCategory     string          `json:"colorCategory,omitempty"`
	ColorType         string          `json:"colorType,omitempty"`
	ColorName         string          `json:"colorName,omitempty"`
	WheelCategory     string          `json:"wheelCategory,omitempty"`
	WheelColorName    string          `json:"wheelColorName,omitempty"`
	InteriorColorName string          `json:"interiorColorName,omitempty"`
	LPValue           string          `json:"lp_value,omitempty"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt
}

type PlayerRaceStats struct {
	Id               string  `json:"id" gorm:"unique;default:uuid_generate_v4();primaryKey,omitempty"`
	PlayerId         string  `json:"playerId,omitempty"`
	Player           Player  `json:"-" gorm:"references:PlayerId;constraint:OnDelete:CASCADE"`
	ArenaId          *string `json:"arenaId,omitempty"`
	Arena            Arena   `json:"-" gorm:"references:ArenaId;constraint:OnDelete:CASCADE"`
	WinStreak        int64   `json:"winStreak"`
	LoseStreak       int64   `json:"loseStreak"`
	DistanceTraveled float64 `json:"distanceTraveled"`
	ShdWon           int64   `json:"showDownWon"`
	TotalShdPlayed   int64   `json:"totalShdPlayed"`
	TdWon            int64   `json:"takeDownWon"`
	TotalTdPlayed    int64   `json:"totalTdPlayed"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt
}
