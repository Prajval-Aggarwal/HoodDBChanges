package model

type PlayerLevel struct {
	Id         string `json:"id" gorm:"unique;default:uuid_generate_v4();primaryKey,omitempty"`
	Level      int64  `json:"level"`
	XPRequired int64  `json:"xpRequired"`
	Coins      int64  `json:"coins"`
}
