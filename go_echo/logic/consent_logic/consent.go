package consent_logic

import (
	"context"
	"database/sql"
	"github.com/tsrkzy/jump_in/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func InsertConsent(ctx context.Context, tx *sql.Tx, eId int64, adminAccountId int64, eventAccountId int64, message string) error {
	consent := models.Consent{
		EventID:        eId,
		AdminAccountID: adminAccountId,
		AccountID:      eventAccountId,
		Message:        message,
		Accepted:       false,
	}

	err := consent.Insert(ctx, tx, boil.Infer())
	if err != nil {
		return err
	}
	return nil
}

func UpdateConsentAccepted(ctx context.Context, tx *sql.Tx, consent *models.Consent) error {
	consent.Accepted = true
	_, err := consent.Update(ctx, tx, boil.Infer())
	return err
}

func FetchConsentByUK(ctx context.Context, tx *sql.Tx, consentId int64, accountId int64, eventId int64) (*models.Consent, error) {
	consent, err := models.Consents(qm.Where("id = ? and account_id = ? and event_id = ?",
		consentId,
		accountId,
		eventId)).One(ctx, tx)
	return consent, err
}
