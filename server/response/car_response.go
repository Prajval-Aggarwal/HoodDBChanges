package response

type CarRes struct {
	CustId   string `json:"custId,omitempty"`
	CarId    string `json:"carId,omitempty"`
	CarName  string `json:"carName"`
	Rarity   int64  `json:"rarity,omitempty"`
	Defaults struct {
		Stats    Stats `json:"stats,omitempty"`
		Purchase struct {
			CurrencyType int64 `json:"currencyType,omitempty"` // 1 for coins and 2 fro cash
			Amount       int64 `json:"amount,omitempty"`
			PremiumBuy   int64 `json:"premiumBuy,omitempty"`
		} `json:"price,omitempty"`
		Customization Customization `json:"carLooks"`
	} `json:"defaultData,omitempty"`
	CarCurrentData struct {
		Stats         Stats         `json:"stats,omitempty"`
		Customization Customization `json:"carLooks,omitempty"`
	} `json:"currentData,omitempty"`
	Status struct {
		Purchasable bool `json:"purchasable"`
		Owned       bool `json:"owned"`
	} `json:"status,omitempty"`
}

type Stats struct {
	Power       int64   "json:\"power,omitempty\""
	Grip        int64   "json:\"grip,omitempty\""
	Weight      int64   "json:\"weight,omitempty\""
	ShiftTime   float64 "json:\"shiftTime,omitempty\""
	OVR         float64 "json:\"ovr,omitempty\""
	Durability  int64   "json:\"durability,omitempty\""
	NitrousTime float64 "json:\"nitrousTime,omitempty\""
}

type Customization struct {
	ColorCategory     string `json:"colorCategory,omitempty"`
	ColorType         int64  `json:"paintType,omitempty"`
	ColorName         int64  `json:"colorType,omitempty"`
	WheelCategory     string `json:"wheelCategory,omitempty"`
	WheelColorName    int64  `json:"caliperType,omitempty"`
	InteriorColorName int64  `json:"interiorType,omitempty"`
	LPValue           string `json:"lp_value"`
}

type UpgradeResponse struct {
	CarLevel      int64 `json:"carLevel"`
	NextLevelCost int64 `json:"nextLevelCost"`
	NewStats      Stats `json:"newStats"`
	NextStats     Stats `json:"nextStats"`
	IsUpgradable  bool  `json:"isUpgradable"`
}

type CarCustom struct {
	CustId      string  `json:"custId"  gorm:"unique;default:uuid_generate_v4();primaryKey,omitempty"`
	CarId       string  `json:"carId,omitempty"`
	Power       int64   `json:"power,omitempty"`
	Grip        int64   `json:"grip,omitempty"`
	ShiftTime   float64 `json:"shiftTime,omitempty"`
	Weight      int64   `json:"weight,omitempty"`
	OVR         float64 `json:"or,omitempty"` //overall rating of the car
	Durability  int64   `json:"Durability,omitempty"`
	NitrousTime float64 `json:"nitrousTime,omitempty"` //increased when nitrous is upgraded
	Class       int64   `json:"class"`
	CarName     string  `json:"carName"`
}
