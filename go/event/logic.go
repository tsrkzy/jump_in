package event

import (
	"context"
	"database/sql"
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
		ee := Event{*e}
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
		ee := Event{*e}
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
		ee := Event{*e}
		eventList = append(eventList, ee)
	}

	return eventList, nil
}
