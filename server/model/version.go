package model

import "time"

type DbVersion struct {
	Version   int `json:"version"`
	UpdatedAt time.Time
}
