package model

import "gorm.io/gorm"

type PlayerLevel struct {
	gorm.Model
	Level      int64 `json:"level"`
	XPRequired int64 `json:"xpRequired"`
	Coins      int64 `json:"coins"`
}
