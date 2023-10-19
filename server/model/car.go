package model

type Car struct {
	CarId      string  `json:"carId"  gorm:"unique;primaryKey"`
	CarName    string  `json:"carName,omitempty"`
	CurrType   string  `json:"currType,omitempty" `
	CurrAmount float64 `json:"cost,omitempty"`
	PremiumBuy int64   `json:"premiumBuy,omitempty"`
	Class      int64   `json:"class,omitempty"`
	Locked     bool    `json:"locked,omitempty"`
}
type DefaultCustomisation struct {
	Id                string  `json:"id" gorm:"unique;default:uuid_generate_v4();primaryKey,omitempty"`
	CarId             string  `json:"carId,omitempty"`
	Car               Car     `json:"-" gorm:"references:CarId;constraint:OnDelete:CASCADE"`
	Power             int64   `json:"power,omitempty"`
	Grip              int64   `json:"grip,omitempty"`
	ShiftTime         float64 `json:"shiftTime,omitempty"`
	Weight            int64   `json:"weight,omitempty"`
	OVR               float64 `json:"or,omitempty"` //overall rating of the car
	Durability        int64   `json:"Durability,omitempty"`
	NitrousTime       float64 `json:"nitrousTime,omitempty"` //increased when nitrous is upgraded
	ColorCategory     string  `json:"colorCategory,omitempty"`
	ColorType         string  `json:"colorType,omitempty"`
	ColorName         string  `json:"colorName,omitempty"`
	WheelCategory     string  `json:"wheelCategory,omitempty"`
	WheelColorName    string  `json:"wheelColorName,omitempty"`
	InteriorColorName string  `json:"interiorColorName,omitempty"`
	LPValue           string  `json:"lpValue,omitempty"`
}

type PartCustomization struct {
	Id                string  `json:"id" gorm:"unique;default:uuid_generate_v4();primaryKey,omitempty"`
	ColorCategory     string  `json:"colorCategory,omitempty"`
	ColorType         string  `json:"colorType,omitempty"`
	ColorName         string  `json:"colorName,omitempty"`
	WheelCategory     string  `json:"wheelCategory,omitempty"`
	WheelColorName    string  `json:"wheelColorName,omitempty"`
	InteriorColorName string  `json:"interiorColorName,omitempty"`
	LPValue           string  `json:"lpValue,omitempty"`
	CurrType          string  `json:"currType,omitempty"`
	CurrAmount        float64 `json:"currAmount,omitempty"`
}
