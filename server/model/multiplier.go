package model

type RatingMulti struct {
	Id           string  `json:"id" gorm:"unique;default:uuid_generate_v4();primaryKey,omitempty"`
	Class        string  `json:"class"`
	ORMultiplier float64 `json:"orMultiplier"`
	
}
