package model

// DB model to store session information
type Session struct {
	SessionId   string `json:"sessionId" gorm:"default:uuid_generate_v4()"`
	PlayerId    string `json:"userId" `
	SessionType int64  `json:"sessionType"`
	DeviceId    string `json:"deviceId"`
	Token       string `json:"token"`
}

type ResetSession struct {
	Id    string `json:"id"`
	Token string `json:"token"`
}
