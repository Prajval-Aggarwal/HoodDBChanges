package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type AddArenaRequest struct {
	ArenaName  string  `json:"arenaName,omitempty"`
	ArenaLevel int64   `json:"arenaLevel,omitempty"`
	Latitude   float64 `json:"latitude,omitempty"`
	Longitude  float64 `json:"longitude,omitempty"`
}

func (a AddArenaRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.ArenaName, validation.Required),
		validation.Field(&a.ArenaLevel, validation.Required, validation.Min(1), validation.Max(3)),
		// Validate Latitude: must be between -90 and 90 degrees
		validation.Field(&a.Latitude, validation.Required, validation.Min(-90.0), validation.Max(90.0)),
		// Validate Longitude: must be between -180 and 180 degrees
		validation.Field(&a.Longitude, validation.Required, validation.Min(-180.0), validation.Max(180.0)),
	)
}

type DeletArenaReq struct {
	ArenaId string `json:"arenaId"`
}

func (a DeletArenaReq) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.ArenaId, validation.Required),
	)
}

type UpdateArenaReq struct {
	ArenaId    string  `json:"arenaId"`
	ArenaName  string  `json:"arenaName,omitempty"`
	ArenaLevel int64   `json:"arenaLevel,omitempty"`
	Latitude   float64 `json:"latitude,omitempty"`
	Longitude  float64 `json:"longitude,omitempty"`
}

func (a UpdateArenaReq) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.ArenaId, validation.Required),
		// Validate Latitude: must be between -90 and 90 degrees
		validation.Field(&a.Latitude, validation.Min(-90.0), validation.Max(90.0)),
		// Validate Longitude: must be between -180 and 180 degrees
		validation.Field(&a.Longitude, validation.Min(-180.0), validation.Max(180.0)),
	)
}

type ReplaceReq struct {
	ArenaId        string `json:"arenaId"`
	NewCustId      string `json:"newCustId"`
	ExistingCustId string `json:"existingCustId"`
}

func (a ReplaceReq) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.ArenaId, validation.Required),
		validation.Field(&a.NewCustId, validation.Required),
		validation.Field(&a.ExistingCustId, validation.Required),
	)
}

type AddCarArenaRequest struct {
	ArenaId string `json:"arenaId"`
	CustId  string `json:"custId"`
}

func (a AddCarArenaRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.ArenaId, validation.Required),
		validation.Field(&a.CustId, validation.Required),
	)
}

type GetArenaReq struct {
	ArenaId string `json:"arenaId"`
}

func (a GetArenaReq) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.ArenaId, validation.Required))
}

type EndChallengeReq struct {
	ArenaId   string `json:"arenaId"`
	PlayerId1 string `json:"playerId"`
	CustId    string `json:"custId"`
	Seconds   int64  `json:"seconds"`
	MilliSec  int64  `json:"milliSec"`
	MicroSec  int64  `json:"microSec"`
	RaceId    string `json:"raceId"`
}

func (a EndChallengeReq) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.CustId, validation.Required),
		validation.Field(&a.PlayerId1, validation.Required),
		validation.Field(&a.Seconds, validation.Required),
		validation.Field(&a.MilliSec, validation.Required),
		validation.Field(&a.MicroSec, validation.Required),
		validation.Field(&a.RaceId, validation.Required),
	)
}
