package entity

import (
	"fmt"
	"github.com/tsrkzy/jump_in/models"
)

type Account struct {
	ID string `json:"id"`
	models.Account
}

// CreateAccount
// IDカラムをint64からstringへ倒す処理を含むファクトリ
func CreateAccount(a *models.Account) *Account {
	account := Account{Account: *a}
	account.ID = fmt.Sprintf("%d", account.Account.ID)
	return &account
}

type Attend struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	EventID   string `json:"event_id"`
	models.Attend
}

func CreateAttend(a *models.Attend) *Attend {
	attend := Attend{Attend: *a}
	attend.ID = fmt.Sprintf("%d", attend.Attend.ID)
	attend.AccountID = fmt.Sprintf("%d", attend.Attend.AccountID)
	attend.EventID = fmt.Sprintf("%d", attend.Attend.EventID)
	return &attend
}

type MailAccount struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	models.MailAccount
}

// CreateMailAccount
// IDカラムをint64からstringへ倒す処理を含むファクトリ
func CreateMailAccount(ma *models.MailAccount) *MailAccount {
	mailAccount := MailAccount{MailAccount: *ma}
	mailAccount.ID = fmt.Sprintf("%d", mailAccount.MailAccount.ID)
	mailAccount.AccountID = fmt.Sprintf("%d", mailAccount.MailAccount.AccountID)

	return &mailAccount
}

type Candidate struct {
	ID      string `json:"id"`
	EventID string `json:"event_id"`
	Votes   []Vote `json:"votes"`
	models.Candidate
}

// CreateCandidate
// IDカラムをint64からstringへ倒す処理を含むファクトリ
func CreateCandidate(c *models.Candidate) *Candidate {
	candidate := Candidate{Candidate: *c}
	candidate.ID = fmt.Sprintf("%d", candidate.Candidate.ID)
	candidate.EventID = fmt.Sprintf("%d", candidate.Candidate.EventID)
	candidate.Votes = make([]Vote, 0)
	return &candidate
}

type Vote struct {
	ID          string  `json:"id"`
	AccountID   string  `json:"account_id"`
	CandidateID string  `json:"candidate_id"`
	Account     Account `json:"account"`
	models.Vote
}

func CreateVote(v *models.Vote) *Vote {
	vote := Vote{Vote: *v}
	vote.ID = fmt.Sprintf("%d", vote.Vote.ID)
	vote.AccountID = fmt.Sprintf("%d", vote.Vote.AccountID)
	vote.CandidateID = fmt.Sprintf("%d", vote.Vote.CandidateID)
	return &vote
}

type Event struct {
	ID           string `json:"id"`
	AccountID    string `json:"account_id"`
	EventGroupID string `json:"event_group_id"`
	models.Event
}

// CreateEvent
// IDカラムをint64からstringへ倒す処理を含むファクトリ
func CreateEvent(e *models.Event) *Event {
	event := Event{Event: *e}
	event.ID = fmt.Sprintf("%d", event.Event.ID)
	event.AccountID = fmt.Sprintf("%d", event.Event.AccountID)
	event.EventGroupID = fmt.Sprintf("%d", event.Event.EventGroupID)
	return &event
}

type Participants struct {
	Attend  Attend  `json:"attend"`
	Account Account `json:"account"`
}

func CreateParticipants(a *Account, att *Attend) *Participants {
	p := &Participants{
		Attend:  *att,
		Account: *a,
	}

	return p
}

type EventDetail struct {
	Event
	Candidates   []Candidate    `json:"candidates"`
	Owner        Account        `json:"owner"`
	Participants []Participants `json:"participants"`
}
