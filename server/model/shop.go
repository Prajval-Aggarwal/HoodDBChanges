package model

import "gorm.io/gorm"

type Shop struct {
	gorm.Model
	PurchaseType  int64 `json:"purchaseType"`
	RewardType    int64 `json:"rewardType"`
	PurchaseValue int64 `json:"purchaseValue"`
	RewardValue   int64 `json:"rewardValue"`
}
