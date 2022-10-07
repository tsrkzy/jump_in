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
