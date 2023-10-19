package model

import "gorm.io/gorm"

type RatingMulti struct {
	gorm.Model
	Class        string  `json:"class"`
	ORMultiplier float64 `json:"orMultiplier"`
}
