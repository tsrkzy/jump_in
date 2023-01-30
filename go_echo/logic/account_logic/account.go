package account_logic

import (
	"context"
	"database/sql"
	"github.com/tsrkzy/jump_in/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func UpdateAccount(err error, a *models.Account, ctx context.Context, tx *sql.Tx) error {
	_, err = a.Update(ctx, tx, boil.Infer())
	return err
}
