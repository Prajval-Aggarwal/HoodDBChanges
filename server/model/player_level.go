package model

import (
	"time"

	"gorm.io/gorm"
)

type PlayerLevel struct {
	Id         string `json:"id" gorm:"unique;default:uuid_generate_v4();primaryKey,omitempty"`
	Level      int64  `json:"level"`
	XPRequired int64  `json:"xpRequired"`
	Coins      int64  `json:"coins"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}
