package helper

import (
	"fmt"
	"github.com/tsrkzy/jump_in/logic/lg"
	"github.com/tsrkzy/jump_in/types/entity"
	"strings"
)

type Maskable interface {
	GetMailAccounts() []entity.MailAccount
}

// Mask (Maskable)
// センシティブな情報(メールアドレス)をマスクする
func Mask(m Maskable) {
	mailAccounts := m.GetMailAccounts()
	for i := 0; i < len(mailAccounts); i++ {
		ma := &(mailAccounts[i])
		mailAddress := ma.MailAddress
		ma.MailAddress = MaskMailAddress(mailAddress)
		lg.Debug(ma.MailAddress)
	}
}

// MaskMailAddress
// メールアドレスのマスク用
// アカウント(@以前)を先頭と末尾のみ残してアスタリスクで伏せる
func MaskMailAddress(s string) string {
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
