package model

///this model or table is created as in player table we are not storing the password

type Admin struct {
	Id       string `json:"id" gorm:"default:uuid_generate_v4();primaryKey"`
	Username string `json:"username"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"`
}
