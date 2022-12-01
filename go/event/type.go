package event

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/tsrkzy/jump_in/models"
)

type ListRequest struct{}

type Event struct {
	models.Event
}

type ListResponse struct {
	EventsOwns    []Event `json:"events_owns"`
	EventsJoins   []Event `json:"events_joins"`
	EventsRunning []Event `json:"events_running"`
}

type CreateRequest struct {
	Name string
}

func (r CreateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.Name,
			validation.Required.Error("イベント名は必須です"),
			validation.RuneLength(5, 40).Error("イベント名は5〜40文字で指定してください"),
		),
	)
}

type CreateResponse struct {
	Event
}

type DetailRequest struct {
	EventId string `query:"event_id"`
}

func (r DetailRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.EventId,
			validation.Required.Error("イベントIDは必須です"),
			is.Int,
		),
	)
}

type User struct {
	models.Account
}

type DetailResponse struct {
	Event
	Owner   User   `json:"owner"`
	Members []User `json:"members"`
}

type AttendRequest struct {
	EventId string `json:"event_id"`
}

func (r AttendRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.EventId,
			validation.Required.Error("イベントIDは必須です"),
			is.Int,
		),
	)
}

type LeaveRequest struct {
	EventId string `json:"event_id"`
}

func (r LeaveRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.EventId,
			validation.Required.Error("イベントIDは必須です"),
			is.Int,
		),
	)
}
