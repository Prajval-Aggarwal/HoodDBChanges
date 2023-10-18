package request

import validation "github.com/go-ozzo/ozzo-validation"

type LicenseRequest struct {
	CarId string `json:"carId"`
	Value string `json:"value"`
}
type InteriorCustomizeRequest struct {
	CarId     string `json:"carId"`
	ColorName string `json:"colorName"`
	ColorCode string `json:"colorCode"`
}

type WheelCustomizeRequest struct {
	CarId         string `json:"carId"`
	ColorCategory string `json:"colorCategory"`
	ColorName     string `json:"colorName"`
	ColorCode     string `json:"colorCode"`
}

type ColorCustomizationRequest struct {
	CarId         string `json:"carId"`
	ColorCategory string `json:"colorCategory"`
	ColorType     string `json:"colorType"`
	ColorName     string `json:"colorName"`
	ColorCode     string `json:"colorCode"`
}

type GetCarColorCategoriesRequest struct {
	Part string `json:"part"`
}
type GetCarColorTypesRequest struct {
	Part          string `json:"part"`
	ColorCategory string `json:"colorCategory"`
}

type GetCarColorRequest struct {
	Part          string `json:"part"`
	ColorCategory string `json:"colorCategory"`
	ColorType     string `json:"colorType"`
}

func (a LicenseRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.CarId, validation.Required),
		validation.Field(&a.Value, validation.Required),
	)
}

func (a InteriorCustomizeRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.CarId, validation.Required),
		validation.Field(&a.ColorName, validation.Required),
		validation.Field(&a.ColorCode, validation.Required),
	)
}

func (a WheelCustomizeRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.CarId, validation.Required),
		validation.Field(&a.ColorCategory, validation.Required),
		validation.Field(&a.ColorName, validation.Required),
		validation.Field(&a.ColorCode, validation.Required),
	)
}

func (a ColorCustomizationRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.CarId, validation.Required),
		validation.Field(&a.ColorCategory, validation.Required),
		validation.Field(&a.ColorType, validation.Required),
		validation.Field(&a.ColorName, validation.Required),
		validation.Field(&a.ColorCode, validation.Required),
	)
}

func (a GetCarColorCategoriesRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Part, validation.Required),
	)
}

func (a GetCarColorTypesRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Part, validation.Required),
		validation.Field(&a.ColorCategory, validation.Required),
	)
}

func (a GetCarColorRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Part, validation.Required),
		validation.Field(&a.ColorCategory, validation.Required),
		validation.Field(&a.ColorType, validation.Required),
	)
}
