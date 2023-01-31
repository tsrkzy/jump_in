package event_logic

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/tsrkzy/jump_in/logic/lg"
	"github.com/tsrkzy/jump_in/models"
	"github.com/tsrkzy/jump_in/types/entity"
	"github.com/tsrkzy/jump_in/types/response_types"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
)

// createEventGroup event_group が存在しない場合は作成する
// event_group.name で検索し、存在しなければ INSERT
func CreateEventGroup(ctx context.Context, tx *sql.Tx, eName string) (*models.EventGroup, error) {
	eg, err := models.EventGroups(qm.Where("name = ?", eName)).One(ctx, tx)

	if err != nil {
		eg = &models.EventGroup{Name: eName}

		err = eg.Insert(ctx, tx, boil.Infer())
		if err != nil {
			return nil, err
		}
	} else {
		lg.Debug("event_group already exists")
	}

	return eg, nil
}

func FetchEventOwns(ctx context.Context, tx *sql.Tx, accountId int64, limit int) ([]entity.Event, error) {
	eOwns, err := models.Events(
		qm.Where("account_id = ?", accountId),
		qm.Limit(limit),
		qm.OrderBy("created_at desc"),
	).All(ctx, tx)
	if err != nil {
		return make([]entity.Event, 0), err
	}

	eventOwns := make([]entity.Event, 0)
	for _, e := range eOwns {
		ee := *entity.CreateEvent(e)
		eventOwns = append(eventOwns, ee)
	}

	return eventOwns, nil
}

func FetchEventJoins(ctx context.Context, tx *sql.Tx, accountId int64, limit int) ([]entity.Event, error) {
	// attend
	at, err := models.Attends(
		qm.Where("account_id = ?", accountId),
		qm.Limit(limit),
		qm.OrderBy("created_at desc"),
	).All(ctx, tx)
	if err != nil {
		return make([]entity.Event, 0), err
	}

	eIdList := make([]interface{}, 0)
	for _, a := range at {
		eId := a.EventID
		eIdList = append(eIdList, int(eId))
	}

	eJoins, err := models.Events(
		qm.WhereIn("id in ?", eIdList...),
	).All(ctx, tx)
	if err != nil {
		return make([]entity.Event, 0), err
	}

	eventJoins := make([]entity.Event, 0)
	for _, e := range eJoins {
		ee := *entity.CreateEvent(e)
		eventJoins = append(eventJoins, ee)
	}

	return eventJoins, nil
}

func FetchNewEventsWithout(ctx context.Context, tx *sql.Tx, eIdExclude []interface{}) ([]entity.Event, error) {
	eList, err := models.Events(
		qm.WhereNotIn("id not in ?", eIdExclude...),
		qm.Where("is_open = ?", true),
		qm.Limit(10),
		qm.OrderBy("created_at DESC"),
	).All(ctx, tx)
	if err != nil {
		return make([]entity.Event, 0), err
	}

	eventList := make([]entity.Event, 0)
	for _, e := range eList {
		ee := *entity.CreateEvent(e)
		eventList = append(eventList, ee)
	}

	return eventList, nil
}

func getOwner(ctx context.Context, tx *sql.Tx, eventId string) (entity.Account, error) {
	e, err := models.Events(qm.Where("id = ?", eventId)).One(ctx, tx)
	if err != nil {
		return entity.Account{}, err
	}
	ownerId := e.AccountID
	aIdList := []int64{ownerId}
	users, err := getUsers(ctx, tx, aIdList)
	if err != nil {
		return entity.Account{}, err
	}
	owner := users[0]
	return owner, nil
}

func getCandidates(ctx context.Context, tx *sql.Tx, eventId string) ([]entity.Candidate, error) {
	candidates, err := models.Candidates(qm.Where("event_id = ?", eventId)).All(ctx, tx)
	if err != nil {
		return nil, err
	}

	// 候補日ID
	cIdList := make([]interface{}, 0)
	for _, c := range candidates {
		cIdList = append(cIdList, c.ID)
	}

	votes, err := models.Votes(qm.WhereIn("candidate_id in ?", cIdList...)).All(ctx, tx)
	if err != nil {
		return nil, err
	}

	// voteからアカウントIDを抽出
	aIdList := make([]interface{}, 0)
	aIdMap := make(map[int64]bool)
	for _, v := range votes {
		if aIdMap[v.AccountID] {
			continue
		}
		aIdMap[v.AccountID] = true
		aIdList = append(aIdList, v.AccountID)
	}

	accounts, err := models.Accounts(qm.WhereIn("id in ?", aIdList...)).All(ctx, tx)
	aMap := make(map[int64]*entity.Account)
	for _, a := range accounts {
		aMap[a.ID] = entity.CreateAccount(a)
	}

	// modelをentityへmap
	cList := make([]entity.Candidate, 0)
	for _, c := range candidates {
		cList = append(cList, *entity.CreateCandidate(c))
	}

	for _, _v := range votes {

		v := entity.CreateVote(_v)
		v.Account = *aMap[v.Vote.AccountID]

		for i := range cList {
			c := &cList[i]
			if v.CandidateID != c.ID {
				continue
			}
			c.Votes = append(c.Votes, *v)
			break
		}
	}
	return cList, nil
}

func getParticipants(ctx context.Context, tx *sql.Tx, eventId string) ([]entity.Participants, error) {
	attendList, err := models.Attends(qm.Where("event_id = ?", eventId)).All(ctx, tx)
	if err != nil {
		return nil, err
	}

	accountIdList := make([]int64, 0)
	for _, a := range attendList {
		accountIdList = append(accountIdList, a.AccountID)
	}

	users, err := getUsers(ctx, tx, accountIdList)
	if err != nil {
		return nil, err
	}

	participants := accountsToParticipants(users, attendList)

	return participants, nil
}

func accountsToParticipants(accounts []entity.Account, attends models.AttendSlice) []entity.Participants {

	attendMap := make(map[int64]*entity.Attend, 0)
	for _, a := range attends {
		attendMap[a.AccountID] = entity.CreateAttend(a)
	}

	participants := make([]entity.Participants, 0)
	for i := range accounts {
		account := &accounts[i]
		att := attendMap[account.Account.ID]
		p := entity.CreateParticipants(account, att)
		participants = append(participants, *p)
	}

	return participants
}

func getUsers(ctx context.Context, tx *sql.Tx, accountIdList []int64) ([]entity.Account, error) {
	aIn := make([]interface{}, 0)
	for _, a := range accountIdList {
		aIn = append(aIn, (interface{})(a))
	}

	aList, err := models.Accounts(qm.WhereIn("id in ?", aIn...)).All(ctx, tx)
	if err != nil {
		return nil, err
	}

	users := make([]entity.Account, 0)
	for _, a := range aList {
		u := *entity.CreateAccount(a)
		users = append(users, u)
	}

	return users, nil
}

func GetDetail(ctx context.Context, tx *sql.Tx, eventId string) (*entity.EventDetail, error) {
	e, err := models.Events(qm.Where("id = ?", eventId)).One(ctx, tx)
	if err != nil {
		msg := fmt.Sprintf("イベントが見つかりません: %s", eventId)
		return nil, response_types.NewErrorSeed(http.StatusNotFound, msg)
	}

	owner, err := getOwner(ctx, tx, eventId)
	if err != nil {
		return nil, err
	}

	candidates, err := getCandidates(ctx, tx, eventId)
	if err != nil {
		return nil, err
	}

	participants, err := getParticipants(ctx, tx, eventId)
	if err != nil {
		return nil, err
	}

	ee := entity.CreateEvent(e)

	dr := &entity.EventDetail{
		Event:        *ee,
		Candidates:   candidates,
		Owner:        owner,
		Participants: participants,
	}

	return dr, nil
}

func InsertEvent(ctx context.Context, tx *sql.Tx, eName string, eDescription string, eg *models.EventGroup, a *models.Account) (*models.Event, error) {
	e := &models.Event{
		Name:         eName,
		Description:  eDescription,
		EventGroupID: eg.ID,
		AccountID:    a.ID,
	}
	err := e.Insert(ctx, tx, boil.Infer())
	return e, err
}

func FetchEventByID(ctx context.Context, tx *sql.Tx, eId int64) (*models.Event, error) {
	e, err := models.Events(qm.Where("id = ?", eId)).One(ctx, tx)
	return e, err
}

func UpsertAttend(attend *models.Attend, ctx context.Context, tx *sql.Tx) error {
	// @REF https://pkg.go.dev/github.com/volatiletech/sqlboiler/v4#readme-upsert
	err := attend.Upsert(ctx, tx, true, []string{"account_id", "event_id"}, boil.Infer(), boil.Infer())
	return err
}

func CreateAttend(ctx context.Context, tx *sql.Tx, aId int64, eId int64, comment string) (*models.Attend, error) {
	att := models.Attend{
		AccountID: aId,
		EventID:   eId,
		Comment:   comment,
	}

	err := att.Insert(ctx, tx, boil.Infer())
	return &att, err
}

func FetchAttendByID(ctx context.Context, tx *sql.Tx, aId int64, eId int64) (*models.Attend, error) {
	att, err := models.Attends(qm.Where("event_id = ? and account_id = ?", eId, aId)).One(ctx, tx)
	return att, err
}

func AttendExists(ctx context.Context, tx *sql.Tx, aId int64, eId int64) (bool, error) {
	exA, err := models.Attends(qm.Where("event_id = ? and account_id = ?", eId, aId)).Exists(ctx, tx)
	return exA, err
}

func DeleteAttendByIDs(ctx context.Context, tx *sql.Tx, aId int64, eId int64) error {
	_, err := models.Attends(qm.Where("event_id = ? and account_id = ?", eId, aId)).DeleteAll(ctx, tx)
	return err
}

func UpdateEvent(ctx context.Context, tx *sql.Tx, e *models.Event) error {
	_, err := e.Update(ctx, tx, boil.Infer())
	return err
}

func CertifyEvent(ctx context.Context, tx *sql.Tx, eventId int64, certify bool) (int64, error) {
	/* 認可予定の event を取得 */
	e, err := FetchEventByID(ctx, tx, eventId)
	if err != nil {
		return 0, response_types.NewErrorSeed(http.StatusNotFound, fmt.Sprintf("イベントが存在しません: %d", eventId))
	}
	e.Certified = certify

	return e.Update(ctx, tx, boil.Infer())
}
