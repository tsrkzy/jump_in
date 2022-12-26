package event

import (
	"context"
	"database/sql"
	"github.com/tsrkzy/jump_in/authenticate"
	"github.com/tsrkzy/jump_in/lg"
	"github.com/tsrkzy/jump_in/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// createEventGroup event_group が存在しない場合は作成する
// event_group.name で検索し、存在しなければ INSERT
func createEventGroup(ctx context.Context, tx *sql.Tx, eName string) (*models.EventGroup, error) {
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

func fetchEventOwns(ctx context.Context, tx *sql.Tx, accountId int64, limit int) ([]Event, error) {
	eOwns, err := models.Events(
		qm.Where("account_id = ?", accountId),
		qm.Limit(limit),
		qm.OrderBy("created_at desc"),
	).All(ctx, tx)
	if err != nil {
		return make([]Event, 0), err
	}

	eventOwns := make([]Event, 0)
	for _, e := range eOwns {
		ee := *CreateEvent(e)
		eventOwns = append(eventOwns, ee)
	}

	return eventOwns, nil
}

func fetchEventJoins(ctx context.Context, tx *sql.Tx, accountId int64, limit int) ([]Event, error) {
	// attend
	at, err := models.Attends(
		qm.Where("account_id = ?", accountId),
		qm.Limit(limit),
		qm.OrderBy("created_at desc"),
	).All(ctx, tx)
	if err != nil {
		return make([]Event, 0), err
	}

	eIdList := make([]interface{}, 0)
	for _, a := range at {
		eId := a.EventID
		eIdList = append(eIdList, int(eId))
	}

	eJoins, err := models.Events(
		qm.WhereIn("id in ?", eIdList...),
		qm.Limit(3),
	).All(ctx, tx)
	if err != nil {
		return make([]Event, 0), err
	}

	eventJoins := make([]Event, 0)
	for _, e := range eJoins {
		ee := *CreateEvent(e)
		eventJoins = append(eventJoins, ee)
	}

	return eventJoins, nil
}

func fetchNewEventsWithout(ctx context.Context, tx *sql.Tx, eIdExclude []interface{}) ([]Event, error) {
	eList, err := models.Events(
		qm.WhereNotIn("id not in ?", eIdExclude...),
		qm.Limit(10),
		qm.OrderBy("created_at DESC"),
	).All(ctx, tx)
	if err != nil {
		return make([]Event, 0), err
	}

	eventList := make([]Event, 0)
	for _, e := range eList {
		ee := *CreateEvent(e)
		eventList = append(eventList, ee)
	}

	return eventList, nil
}

func getOwner(ctx context.Context, tx *sql.Tx, eventId string) (authenticate.Account, error) {
	e, err := models.Events(qm.Where("id = ?", eventId)).One(ctx, tx)
	if err != nil {
		return authenticate.Account{}, err
	}
	ownerId := e.AccountID
	aIdList := []int64{ownerId}
	users, err := getUsers(ctx, tx, aIdList)
	if err != nil {
		return authenticate.Account{}, err
	}
	owner := users[0]
	return owner, nil
}

func getCandidates(ctx context.Context, tx *sql.Tx, eventId string) ([]Candidate, error) {
	candidates, err := models.Candidates(qm.Where("event_id = ?", eventId)).All(ctx, tx)
	if err != nil {
		return nil, err
	}

	cList := make([]Candidate, 0)
	for _, c := range candidates {
		cList = append(cList, *CreateCandidate(c))
	}

	return cList, nil
}

func getParticipants(ctx context.Context, tx *sql.Tx, eventId string) ([]authenticate.Account, error) {
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

	return users, nil
}

func getUsers(ctx context.Context, tx *sql.Tx, accountIdList []int64) ([]authenticate.Account, error) {
	aIn := make([]interface{}, 0)
	for _, a := range accountIdList {
		aIn = append(aIn, (interface{})(a))
	}

	aList, err := models.Accounts(qm.WhereIn("id in ?", aIn...)).All(ctx, tx)
	if err != nil {
		return nil, err
	}

	users := make([]authenticate.Account, 0)
	for _, a := range aList {
		u := *authenticate.CreateAccount(a)
		users = append(users, u)
	}

	return users, nil
}
