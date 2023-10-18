package request

import validation "github.com/go-ozzo/ozzo-validation"

type GarageRequest struct {
	GarageId string `json:"garageId"`
}

func (a GarageRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.GarageId, validation.Required),
	)
}

// add car to garage request model
type AddCarRequest struct {
	GarageId string `json:"garageId"`
	CarId    string `json:"carId"`
}

func (a AddCarRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.GarageId, validation.Required),
		validation.Field(&a.CarId, validation.Required),
	)
}

type AddGarageRequest struct {
	GarageName    string  `json:"garageName,omitempty"`
	GarageType    int64   `json:"garageType,omitempty"`
	Latitude      float64 `json:"latitude,omitempty"`
	Longitude     float64 `json:"longitude,omitempty"`
	Level         int64   `json:"level,omitempty"` //level required to unlock the garage
	Rarity        int64   `json:"rarity"`
	CoinsRequired int64   `json:"coinsRequired,omitempty"` //coins required to unlock the garage
}

func (a AddGarageRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.GarageName, validation.Required),
		validation.Field(&a.GarageType, validation.Required, validation.Min(1), validation.Max(5)),
		// Validate Latitude: must be between -90 and 90 degrees
		validation.Field(&a.Latitude, validation.Required, validation.Min(-90.0), validation.Max(90.0)),
		// Validate Longitude: must be between -180 and 180 degrees
		validation.Field(&a.Longitude, validation.Required, validation.Min(-180.0), validation.Max(180.0)),
		validation.Field(&a.Level, validation.Required, validation.Min(1), validation.Max(50)),
		validation.Field(&a.CoinsRequired, validation.Required),
		validation.Field(&a.Rarity, validation.Required, validation.Min(1), validation.Max(5)),
	)
}

type DeletGarageReq struct {
	GarageId string `json:"garageId"`
}

func (a DeletGarageReq) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.GarageId, validation.Required),
	)
}

type UpdateGarageReq struct {
	GarageId      string  `json:"garageId"`
	GarageName    string  `json:"garageName,omitempty"`
	GarageType    int64   `json:"garageType"`
	Latitude      float64 `json:"latitude,omitempty"`
	Longitude     float64 `json:"longitude,omitempty"`
	Level         int64   `json:"level,omitempty"` //level required to unlock the garage
	Rarity        int64   `json:"rarity"`
	CoinsRequired int64   `json:"coinsRequired,omitempty"` //coins required to unlock the garage
}

func (a UpdateGarageReq) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.GarageId, validation.Required),
		// Validate Latitude: must be between -90 and 90 degrees
		validation.Field(&a.Latitude, validation.Min(-90.0), validation.Max(90.0)),
		// Validate Longitude: must be between -180 and 180 degrees
		validation.Field(&a.Longitude, validation.Min(-180.0), validation.Max(180.0)),
		validation.Field(&a.Level, validation.Min(1), validation.Max(50)),
	)

}
