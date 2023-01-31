package event_types

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/tsrkzy/jump_in/types/entity"
)

type ListRequest struct {
	AccountID string `query:"account_id"`
}

func (r ListRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.AccountID,
			validation.Required.Error("アカウントIDは必須です"),
		))
}

type ListResponse struct {
	EventsOwns    []entity.Event `json:"events_owns"`
	EventsJoins   []entity.Event `json:"events_joins"`
	EventsRunning []entity.Event `json:"events_running"`
}

type CreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	AccountID   string `json:"account_id"`
}

func (r CreateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.Name,
			validation.Required.Error("イベント名は必須です"),
			validation.RuneLength(5, 40).Error("イベント名は5〜40文字で指定してください"),
		),
		validation.Field(
			&r.Description,
			validation.Required.Error("イベント名は必須です"),
			validation.RuneLength(0, 1000).Error("イベント説明は0〜1000文字で指定してください"),
		),
		validation.Field(
			&r.AccountID,
			validation.Required.Error("アカウントIDは必須です"),
		),
	)
}

type CreateResponse struct {
	entity.Event
}

type UpdateNameRequest struct {
	EventID   string `json:"event_id"`
	EventName string `json:"name"`
}

func (r UpdateNameRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.EventID,
			validation.Required.Error("イベントIDは必須です"),
		), validation.Field(
			&r.EventName,
			validation.Required.Error("イベント名は必須です"),
			validation.RuneLength(5, 40).Error("イベント名は5〜40文字で指定してください"),
		),
	)
}

type UpdateNameResponse struct {
	entity.EventDetail
}

type UpdateDescriptionRequest struct {
	EventID     string `json:"event_id"`
	Description string `json:"description"`
}

func (r UpdateDescriptionRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.EventID,
			validation.Required.Error("イベントIDは必須です"),
		), validation.Field(
			&r.Description,
			validation.Required.Error("イベント説明は必須です"),
			validation.RuneLength(0, 1000).Error("イベント説明は0〜1000文字で指定してください"),
		),
	)
}

type UpdateDescriptionResponse struct {
	entity.EventDetail
}

type UpdateOpenRequest struct {
	EventID string `json:"event_id"`
	IsOpen  bool   `json:"is_open"`
}

func (r UpdateOpenRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.EventID,
			validation.Required.Error("イベントIDは必須です"),
		), validation.Field(
			&r.IsOpen,
			validation.NotNil.Error("イベント募集中フラグは必須です"),
		),
	)
}

type UpdateOpenResponse struct {
	entity.EventDetail
}

type DetailRequest struct {
	EventID string `query:"event_id"`
}

func (r DetailRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.EventID,
			validation.Required.Error("イベントIDは必須です"),
			is.Int,
		),
	)
}

// DetailResponse
type DetailResponse struct {
	entity.EventDetail
}

type AttendRequest struct {
	EventID   string `json:"event_id"`
	AccountID string `json:"account_id"`
	Comment   string `json:"comment"`
}

func (r AttendRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.EventID,
			validation.Required.Error("イベントIDは必須です"),
			is.Int,
		), validation.Field(
			&r.AccountID,
			validation.Required.Error("アカウントIDは必須です"),
		),
	)
}

type AttendResponse struct {
	DetailResponse
}

type LeaveRequest struct {
	EventID   string `json:"event_id"`
	AccountID string `json:"account_id"`
}

func (r LeaveRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.EventID,
			validation.Required.Error("イベントIDは必須です"),
			is.Int,
		), validation.Field(
			&r.AccountID,
			validation.Required.Error("アカウントIDは必須です"),
		),
	)
}

type LeaveResponse struct {
	DetailResponse
}

type CertifyRequest struct {
	EventID   string `json:"event_id"`
	AccountID string `json:"account_id"`
	Certify   bool   `json:"certify"`
}

func (r CertifyRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.EventID,
			validation.Required.Error("イベントIDは必須です"),
			is.Int,
		), validation.Field(
			&r.AccountID,
			validation.Required.Error("アカウントIDは必須です"),
		),
	)
}

type CertifyResponse struct {
	DetailResponse
}
