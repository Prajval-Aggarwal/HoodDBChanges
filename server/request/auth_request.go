package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type GuestLoginRequest struct {
	PlayerName string `json:"playerName"`
	DeviceId   string `json:"deviceId"`
	OS         int64  `json:"os"`
	Token      string `json:"token"`
}
type AdminLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PlayerLoginRequest struct {
	Credential string `json:"credential"`
}

func (a PlayerLoginRequest) Validate() error {

	return validation.ValidateStruct(&a,
		validation.Field(&a.Credential, validation.Required),
	)
}

type UpdateEmailRequest struct {
	Email string `json:"email"`
}
type ForgotPassRequest struct {
	Email string `json:"email" `
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	Password    string `json:"password" `
}

type ResetPasswordRequest struct {
	Password string `json:"password" `
}

// Validation
func (a UpdatePasswordRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.OldPassword, validation.Required),
		validation.Field(&a.Password, validation.Required),
	)
}
func (a ForgotPassRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Email, validation.Required, is.Email),
	)
}

func (a GuestLoginRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.PlayerName, validation.Required),
		validation.Field(&a.DeviceId, validation.Required),
		validation.Field(&a.OS, validation.Required),
	)
}

func (a AdminLoginRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.Password, validation.Required),
	)
}

func (a UpdateEmailRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Email, validation.Required, is.Email),
	)
}

func (a ResetPasswordRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Password, validation.Required),
	)
}
