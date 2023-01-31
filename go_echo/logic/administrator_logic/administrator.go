package administrator_logic

import (
	"context"
	"database/sql"
	"github.com/tsrkzy/jump_in/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func InsertAdministrator(ctx context.Context, tx *sql.Tx, accountId int64, invitationId int64) error {
	admin := &models.Administrator{
		AccountID:    accountId,
		InvitationID: invitationId,
	}
	return admin.Insert(ctx, tx, boil.Infer())
}

func FetchAdministratorByUK(ctx context.Context, tx *sql.Tx, accountId int64, invitationId int64) (*models.Administrator, error) {
	return models.Administrators(qm.Where("account_id = ? and invitation_id = ?", accountId, invitationId)).One(ctx, tx)
}
