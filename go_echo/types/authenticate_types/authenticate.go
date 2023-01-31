package authenticate_types

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/tsrkzy/jump_in/types/entity"
)

type Request struct {
	MailAddress string `json:"mail_address"`
	RedirectURI string `json:"redirect_uri"`
}

func (r Request) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.MailAddress,
			validation.Required.Error("メールアドレスは必須です"),
			validation.RuneLength(5, 40).Error("メールアドレスは5〜40文字で指定してください"),
			is.Email.Error("メールアドレスの書式が不正です"),
		),
		validation.Field(
			&r.RedirectURI,
			is.URL.Error("URIを指定してください"),
		),
	)
}

type Result struct {
	MailAddress string `json:"mail_address"`
	URIHash     string `json:"uri_hash"`
	ChocoChip   string `json:"choco_chip"`
	MagicLink   string `json:"magic_link"`
	IpAddress   string `json:"ip_address"`
}

type WhoAmIResponse struct {
	entity.Account
	Administrator entity.Administrator `json:"admin"`
	MailAccounts  []entity.MailAccount `json:"mail_accounts"`
}

func (w *WhoAmIResponse) GetMailAccounts() []entity.MailAccount {
	return w.MailAccounts
}
