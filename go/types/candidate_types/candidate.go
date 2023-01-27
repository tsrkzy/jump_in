package candidate_types

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/tsrkzy/jump_in/types/entity"
	"regexp"
)

type CreateRequest struct {
	EventID   string `json:"event_id"`
	AccountID string `json:"account_id"`
	OpenAt    string `json:"open_at"`
}

func (r CreateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.EventID,
			validation.Required.Error("イベントIDは必須です"),
		),
		validation.Field(
			&r.AccountID,
			validation.Required.Error("アカウントIDは必須です"),
		),
		validation.Field(
			&r.OpenAt,
			validation.Required.Error("候補日時は必須です"),
			validation.Match(regexp.MustCompile(`^\d{12}$`)).Error("候補日時はYYYYmmddHHii(12文字)で指定してください"),
		),
	)
}

type CreateResponse struct {
	entity.EventDetail
}
type DeleteRequest struct {
	EventID     string `json:"event_id"`
	AccountID   string `json:"account_id"`
	CandidateID string `json:"candidate_id"`
}

func (r DeleteRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.EventID,
			validation.Required.Error("イベントIDは必須です"),
		),
		validation.Field(
			&r.AccountID,
			validation.Required.Error("アカウントIDは必須です"),
		), validation.Field(
			&r.CandidateID,
			validation.Required.Error("候補日IDは必須です"),
		),
	)
}

type DeleteResponse struct {
	entity.EventDetail
}

type UpvoteRequest struct {
	EventID     string `json:"event_id"`
	AccountID   string `json:"account_id"`
	CandidateID string `json:"candidate_id"`
}

func (r UpvoteRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.EventID,
			validation.Required.Error("イベントIDは必須です"),
		),
		validation.Field(
			&r.AccountID,
			validation.Required.Error("アカウントIDは必須です"),
		), validation.Field(
			&r.CandidateID,
			validation.Required.Error("候補日IDは必須です"),
		),
	)
}

type UpvoteResponse struct {
	entity.EventDetail
}

type DownvoteRequest struct {
	EventID     string `json:"event_id"`
	AccountID   string `json:"account_id"`
	CandidateID string `json:"candidate_id"`
}

func (r DownvoteRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.EventID,
			validation.Required.Error("イベントIDは必須です"),
		),
		validation.Field(
			&r.AccountID,
			validation.Required.Error("アカウントIDは必須です"),
		), validation.Field(
			&r.CandidateID,
			validation.Required.Error("候補日IDは必須です"),
		),
	)
}

type DownvoteResponse struct {
	entity.EventDetail
}

type UpdateRequest struct {
	EventID    string             `json:"event_id"`
	Candidates []entity.Candidate `json:"candidates"`
}

func (r UpdateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.EventID,
			validation.Required.Error("イベントIDは必須です"),
		),
	)
}

type UpdateResponse struct {
	entity.EventDetail
}
