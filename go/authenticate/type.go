package authenticate

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/tsrkzy/jump_in/lg"
	"github.com/tsrkzy/jump_in/models"
	"strings"
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
	models.Account
	MailAccounts []models.MailAccount `json:"mail_accounts"`
}

// Mask
// センシティブな情報(メールアドレス)をマスクする
func (w *WhoAmIResponse) Mask() {
	for i := 0; i < len(w.MailAccounts); i++ {
		ma := &(w.MailAccounts[i])
		mailAddress := ma.MailAddress
		ma.MailAddress = maskMailAddress(mailAddress)
		lg.Debug(ma.MailAddress)
	}
}

// maskMailAddress
// メールアドレスのマスク用
// アカウント(@以前)を先頭と末尾のみ残してアスタリスクで伏せる
func maskMailAddress(s string) string {
	splits := strings.Split(s, "@")
	account := splits[0]
	domain := splits[1]

	a := ""
	l := len(account)
	if l <= 2 {
		a = account
	} else {
		start := string(account[0])
		end := string(account[l-1])
		asters := strings.Repeat("*", l-2)
		a = fmt.Sprintf("%s%s%s", start, asters, end)
	}

	return fmt.Sprintf("%s@%s", a, domain)
}
