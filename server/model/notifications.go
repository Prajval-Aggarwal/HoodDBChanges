package model

type Notifications struct {
	MessageId   string `json:"messageId" gorm:"default:uuid_generate_v4();primaryKey"`
	PlayerId    string `json:"playerId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      int64  `json:"status"`
}