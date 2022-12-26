package event

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/tsrkzy/jump_in/authenticate"
	"github.com/tsrkzy/jump_in/models"
)

type ListRequest struct {
	AccountId string `query:"account_id"`
}

func (r ListRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.AccountId,
			validation.Required.Error("アカウントIDは必須です"),
		))
}

type Event struct {
	ID           string `json:"id"`
	AccountID    string `json:"account_id"`
	EventGroupID string `json:"event_group_id"`
	models.Event
}

func CreateEvent(e *models.Event) *Event {
	event := Event{Event: *e}
	event.ID = fmt.Sprintf("%d", event.Event.ID)
	event.AccountID = fmt.Sprintf("%d", event.Event.AccountID)
	event.EventGroupID = fmt.Sprintf("%d", event.Event.EventGroupID)
	return &event
}

type Candidate struct {
	ID      string `json:"id"`
	EventID string `json:"event_id"`
	models.Candidate
}

func CreateCandidate(c *models.Candidate) *Candidate {
	candidate := Candidate{Candidate: *c}
	candidate.ID = fmt.Sprintf("%d", candidate.Candidate.ID)
	candidate.EventID = fmt.Sprintf("%d", candidate.Candidate.EventID)
	return &candidate
}

type ListResponse struct {
	EventsOwns    []Event `json:"events_owns"`
	EventsJoins   []Event `json:"events_joins"`
	EventsRunning []Event `json:"events_running"`
}

type CreateRequest struct {
	Name      string `json:"name"`
	AccountId string `json:"account_id"`
}

func (r CreateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.Name,
			validation.Required.Error("イベント名は必須です"),
			validation.RuneLength(5, 40).Error("イベント名は5〜40文字で指定してください"),
		),
		validation.Field(
			&r.AccountId,
			validation.Required.Error("アカウントIDは必須です"),
		),
	)
}

type CreateResponse struct {
	Event
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

type DetailResponse struct {
	Event
	Candidates   []Candidate            `json:"candidates"`
	Owner        authenticate.Account   `json:"owner"`
	Participants []authenticate.Account `json:"participants"`
}

type AttendRequest struct {
	EventId   string `json:"event_id"`
	AccountId string `json:"account_id"`
}

func (r AttendRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.EventId,
			validation.Required.Error("イベントIDは必須です"),
			is.Int,
		), validation.Field(
			&r.AccountId,
			validation.Required.Error("アカウントIDは必須です"),
		),
	)
}

type AttendResponse struct {
	DetailResponse
}

type LeaveRequest struct {
	EventId   string `json:"event_id"`
	AccountId string `json:"account_id"`
}

func (r LeaveRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.EventId,
			validation.Required.Error("イベントIDは必須です"),
			is.Int,
		), validation.Field(
			&r.AccountId,
			validation.Required.Error("アカウントIDは必須です"),
		),
	)
}

type LeaveResponse struct {
	DetailResponse
}
