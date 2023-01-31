package consent_types

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/tsrkzy/jump_in/types/entity"
)

type CreateRequest struct {
	AccountID string `json:"account_id"`
	EventID   string `json:"event_id"`
	Message   string `json:"message"`
}

func (r CreateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.EventID,
			validation.Required.Error("イベントIDは必須です"),
		), validation.Field(
			&r.AccountID,
			validation.Required.Error("アカウントIDは必須です"),
		), validation.Field(
			&r.Message,
			validation.Required.Error("本文は必須です"),
			validation.RuneLength(0, 2000).Error("本文は0〜2000文字で指定してください"),
		),
	)
}

type CreateResponse struct {
	entity.EventDetail
}

type AcceptRequest struct {
	AccountID string `json:"account_id"`
	EventID   string `json:"event_id"`
	ConsentID string `json:"consent_id"`
}

func (r AcceptRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.EventID,
			validation.Required.Error("イベントIDは必須です"),
		), validation.Field(
			&r.AccountID,
			validation.Required.Error("アカウントIDは必須です"),
		), validation.Field(
			&r.ConsentID,
			validation.Required.Error("同意書IDは必須です"),
		),
	)
}

type AcceptResponse struct {
	entity.EventDetail
}
