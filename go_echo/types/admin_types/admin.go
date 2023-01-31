package admin_types

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/tsrkzy/jump_in/types/response_types"
)

type LoginRequest struct {
	AccountID string `json:"account_id"`
	PassHash  string `json:"pass_hash"`
}

func (r LoginRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.AccountID,
			validation.Required,
		),
		validation.Field(
			&r.PassHash,
			validation.Required,
		))
}

type LoginResponse struct {
	response_types.OK
}

type LogoutRequest struct {
	AccountID string `json:"account_id"`
}

func (r LogoutRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.AccountID,
			validation.Required,
		))
}

type LogoutResponse struct {
	response_types.OK
}
