package model

import (
	"time"

	"gorm.io/gorm"
)

type Garage struct {
	GarageId      string    `json:"garageId"  gorm:"unique;default:uuid_generate_v4();primaryKey"`
	GarageName    string    `json:"garageName,omitempty"`
	GarageType    int64     `json:"garageType,omitempty"`
	Latitude      float64   `json:"latitude,omitempty"`
	Longitude     float64   `json:"longitude,omitempty"`
	Level         int64     `json:"level,omitempty"`         //level reuired to unlock the garage
	CoinsRequired int64     `json:"coinsRequired,omitempty"` //coins required to unlock the garage
	Rarity        int64     `json:"rarity"`
	Capacity      int64     `json:"capacity"`
	Locked        bool      `json:"locked,omitempty"`
	CreatedAt     time.Time `json:"createdAt,omitempty"`

	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type GarageCars struct {
	Id        string `json:"id" gorm:"unique;default:uuid_generate_v4();primaryKey,omitempty"`
	PlayerId  string `json:"playerId"`
	Player    Player `json:"-" gorm:"references:PlayerId;constraint:OnDelete:CASCADE"`
	GarageId  string `json:"garageId"`
	Garage    Garage `json:"-" gorm:"references:GarageId;constraint:OnDelete:CASCADE"`
	CustId    string `json:"custId"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
