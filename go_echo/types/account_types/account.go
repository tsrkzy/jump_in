package account_types

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/tsrkzy/jump_in/types/entity"
)

type UpdateNameRequest struct {
	AccountID string `json:"account_id"`
	Name      string `json:"name"`
}

func (r UpdateNameRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.AccountID,
			validation.Required.Error("アカウントIDは必須です"),
		), validation.Field(
			&r.Name,
			validation.Required.Error("アカウント名は必須です"),
			validation.RuneLength(1, 10).Error("イベント名は1〜10文字で指定してください"),
		))

}

type UpdateNameResponse struct {
	entity.Account
	MailAccounts []entity.MailAccount `json:"mail_accounts"`
}

func (w *UpdateNameResponse) GetMailAccounts() []entity.MailAccount {
	return w.MailAccounts
}
