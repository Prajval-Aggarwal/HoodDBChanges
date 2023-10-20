package model

import (
	"time"

	"gorm.io/gorm"
)

type Shop struct {
	Id            string `json:"id" gorm:"unique;default:uuid_generate_v4();primaryKey,omitempty"`
	PurchaseType  int64  `json:"purchaseType"`
	RewardType    int64  `json:"rewardType"`
	PurchaseValue int64  `json:"purchaseValue"`
	RewardValue   int64  `json:"rewardValue"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}
