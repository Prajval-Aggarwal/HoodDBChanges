package model

import (
	"time"

	"gorm.io/gorm"
)

type ArenaRaceRecord struct {
	Id        string    `json:"id" gorm:"unique;default:uuid_generate_v4();primaryKey,omitempty"`
	PlayerId  string    `json:"playerId"`
	Player    Player    `json:"-" gorm:"references:PlayerId;constraint:OnDelete:CASCADE"`
	ArenaId   string    `json:"arenaId"`
	Arena     Arena     `josn:"-" gorm:"references:ArenaId;constraint:OnDelete:CASCADE"`
	TimeWin   string    `json:"time"`
	Result    string    `json:"result"`
	CarId     string    `json:"carId"`
	Car       Car       `json:"-" gorm:"references:CarId;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `json:"createdAt"`
}

type TempRaceRecords struct {
	Id       string `json:"id" gorm:"unique;default:uuid_generate_v4();primaryKey,omitempty"`
	PlayerId string `json:"playerId"`
	Player   Player `json:"-" gorm:"references:PlayerId;constraint:OnDelete:CASCADE"`
	ArenaId  string `json:"arenaId"`
	Arena    Arena  `josn:"-" gorm:"references:ArenaId;constraint:OnDelete:CASCADE"`
	TimeWin  string `json:"time"`
	Result   string `json:"result"`
	CarId    string `json:"carId"`
	Car      Car    `json:"-" gorm:"references:CarId;constraint:OnDelete:CASCADE"`
}
type ArenaLevelPerks struct {
	Id          string `json:"id" gorm:"unique;default:uuid_generate_v4();primaryKey,omitempty"`
	ArenaLevel  int64  `json:"arenaLevel"`
	Coins       int64  `json:"coins"`
	Cash        int64  `json:"cash"`
	RepairParts int64  `json:"repairParts"`
}

type ArenaCars struct {
	Id        string `json:"id" gorm:"unique;default:uuid_generate_v4();primaryKey,omitempty"`
	PlayerId  string `json:"playerId" `
	Player    Player `json:"-" gorm:"references:PlayerId;constraint:OnDelete:CASCADE"`
	ArenaId   string `json:"arenaId"`
	Arena     Arena  `json:"-" gorm:"references:ArenaId;constraint:OnDelete:CASCADE"`
	CustId    string `json:"custId"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
type Arena struct {
	ArenaId    string    `json:"arenaId" gorm:"unique;default:uuid_generate_v4();primaryKey,omitempty"`
	ArenaName  string    `json:"arenaName"`
	ArenaLevel int64     `json:"arenaLevel"`
	Longitude  float64   `json:"longitude"`
	Latitude   float64   `json:"latitude"`
	CreatedAt  time.Time `json:"createdAt,omitempty"`
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}
