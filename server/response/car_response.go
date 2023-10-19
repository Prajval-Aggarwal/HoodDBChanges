package response

type CarRes struct {
	CarId    string `json:"carId,omitempty"`
	CarName  string `json:"carName,omitempty"`
	Rarity   int64  `json:"rarity,omitempty"`
	Defaults struct {
		Stats    Stats `json:"stats,omitempty"`
		Purchase struct {
			CurrencyType int64 `json:"currencyType,omitempty"` // 1 for coins and 2 fro cash
			Amount       int64 `json:"amount,omitempty"`
			PremiumBuy   int64 `json:"premiumBuy,omitempty"`
		} `json:"price,omitempty"`
		Customization []Customization `json:"carLooks,omitempty"`
	} `json:"defaultData,omitempty"`
	CarCurrentData struct {
		Stats         Stats           `json:"stats,omitempty"`
		Customization []Customization `json:"carLooks,omitempty"`
	} `json:"currentData,omitempty"`
	Status struct {
		Purchasable bool `json:"purchasable"`
		Owned       bool `json:"owned"`
	} `json:"status,omitempty"`
}

type Stats struct {
	Power      int64   "json:\"power,omitempty\""
	Grip       int64   "json:\"grip,omitempty\""
	Weight     int64   "json:\"weight,omitempty\""
	ShiftTime  float64 "json:\"shiftTime,omitempty\""
	OVR        float64 "json:\"ovr,omitempty\""
	Durability int64   "json:\"durability,omitempty\""
}

type Customization struct {
	Part          string `json:"part,omitempty"`
	ColorCategory string `json:"colorCategory,omitempty"`
	ColorType     string `json:"colorType,omitempty"`
	ColorName     string `json:"colorName,omitempty"`
	Value         string `json:"value,omitempty"`
}

type UpgradeResponse struct {
	CarLevel      int64 `json:"carLevel"`
	NextLevelCost int64 `json:"nextLevelCost"`
	NewStats      Stats `json:"newStats"`
	NextStats     Stats `json:"nextStats"`
	IsUpgradable  bool  `json:"isUpgradable"`
}
