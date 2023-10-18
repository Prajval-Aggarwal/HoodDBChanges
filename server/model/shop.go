package model

type Shop struct {
	Id            string `json:"id" gorm:"default:uuid_generate_v4();primaryKey"`
	PurchaseType  int64  `json:"purchaseType"`
	RewardType    int64  `json:"rewardType"`
	PurchaseValue int64  `json:"purchaseValue"`
	RewardValue   int64  `json:"rewardValue"`
}
