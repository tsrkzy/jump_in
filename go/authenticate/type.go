package authenticate

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/tsrkzy/jump_in/models"
)

type Request struct {
	Email       string `json:"email"`
	RedirectURI string `json:"redirect_uri"`
}

func (r Request) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.Email,
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
	Email     string `json:"email"`
	URIHash   string `json:"uri_hash"`
	ChocoChip string `json:"choco_chip"`
	MagicLink string `json:"magic_link"`
	IpAddress string `json:"ip_address"`
}

type WhoAmIResponse struct {
	models.Account
	MailAccounts []models.MailAccount `json:"mail_accounts"`
}
