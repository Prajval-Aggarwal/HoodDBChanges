package model

// DB model to store session information
type Session struct {
	SessionId string `json:"sessionId" gorm:"default:uuid_generate_v4();unique;primaryKey"`
	UserId    string `json:"userId"`
	Token     string `json:"token"`
}

type ResetSession struct {
	Id    string `json:"id"`
	Token string `json:"token"`
}
