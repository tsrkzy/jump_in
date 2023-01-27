package authenticate_logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/tsrkzy/jump_in/logic/lg"
	"github.com/tsrkzy/jump_in/models"
	"github.com/tsrkzy/jump_in/sess"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"math/rand"
	"strings"
	"time"
)

// CreateInvitation
// メールアドレスが mail_account に存在するかチェック
// 存在しないなら、 mail_account と account を作成
//
// 入力したメールアドレスに対応する invitation を作成
func CreateInvitation(ctx *context.Context, tx *sql.Tx, invitation *models.Invitation, mailAddress string) error {

	var a *models.Account
	var ma *models.MailAccount
	var err error

	ma, err = models.MailAccounts(qm.Where("mail_address = ?", mailAddress)).One(*ctx, tx)

	if err != nil {
		lg.Debug("create mail_account and account")
		/* mail_account と account がない場合は作成 */
		a = &models.Account{
			Name: fmt.Sprintf("%s", strings.Split(mailAddress, "@")[0]),
		}
		err := a.Insert(*ctx, tx, boil.Infer())
		if err != nil {
			return err
		}

		ma = &models.MailAccount{
			AccountID:   a.ID,
			MailAddress: mailAddress,
		}
		err = ma.Insert(*ctx, tx, boil.Infer())
		if err != nil {
			return err
		}
		lg.Debugf("account_id: %d", a.ID)
		lg.Debugf("mail_address: %s", ma.MailAddress)
		lg.Debugf("mail_account_id: %d", ma.ID)
	}

	invitation.MailAccountID = ma.AccountID

	err = invitation.Insert(*ctx, tx, boil.Infer())
	return err
}

// GenerateHash
// hyphen-separated 24(6x4) characters
// xxxxxx-xxxxxx-xxxxxx-xxxxxx
// 123456 123456 123456 123456
//
// regexp pattern: /^l[a-z0-9]{5}(\-[a-z0-9]{6}){3}$/
func GenerateHash(ctx *context.Context, tx *sql.Tx, column string) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	const sections = 4
	const sectionSize = 8

	var hash string
	iterCount := 0
	for true {
		/* 試行回数リミット */
		if iterCount > 10 {
			return "", errors.New("exceed iterate limit in generateHash()")
		}

		h := make([]string, sections*sectionSize)
		for i := 0; i < sections; i++ {
			for j := 0; j < sectionSize; j++ {
				u := letters[rand.Intn(len(letters))]
				h[i*sectionSize+j] = string(u)
			}
		}
		s1 := strings.Join(h[0:sectionSize*1], "")
		s2 := strings.Join(h[sectionSize*1:sectionSize*2], "")
		s3 := strings.Join(h[sectionSize*2:sectionSize*3], "")
		s4 := strings.Join(h[sectionSize*3:sectionSize*4], "")
		s := []string{s1, s2, s3, s4}
		hash = strings.Join(s, "-")

		/* hashの存在チェック */
		if isNew, err := IsNewHash(ctx, tx, column, hash); err == nil && isNew {
			break
		}

		iterCount++
	}

	return hash, nil
}

// IsNewHash
// DBのレコードをさらって、まだ存在しないhashだったらtrue
func IsNewHash(ctx *context.Context, tx *sql.Tx, column string, hash string) (bool, error) {
	exists, err := models.Invitations(qm.Where(column+" = ?", hash)).Exists(*ctx, tx)
	if err != nil {
		return false, err
	}
	return !exists, nil
}

func CheckExistence(ctx *context.Context, tx *sql.Tx, uriHash string, chocoChip string) (bool, error) {
	lg.Debugf("CheckExistence: %s", uriHash)
	r, err := models.Invitations(qm.Where("uri_hash = ?", uriHash)).One(*ctx, tx)
	if err != nil || r == nil {
		/* 存在しないMLだった(Oneはレコードが存在しないと error 扱い) */
		lg.Debugf("no record found with %s", uriHash)
		return false, err
	}

	/* cookie(choco_chip) が一致しているか */
	lg.Debug(r.ChocoChip)
	lg.Debug(chocoChip)
	if r.ChocoChip != chocoChip {
		lg.Debug("choco_chip is wrong.")
		return false, nil
	}

	/* 未使用である */
	lg.Debug(r.AuthorisedDatetime)
	if r.Authorised == true {
		lg.Debug("this link already used")
		return false, nil
	}

	/* 現在時刻がMLの受付期限(expired_datetime)を過ぎていないか */
	lg.Debug(r.ExpiredDatetime)
	lg.Debug(time.Now())
	if time.Now().After(r.ExpiredDatetime) {
		lg.Debug("exceed expired datetime")
		return false, nil
	}

	return true, err
}

func AuthoriseMagicLink(ctx *context.Context, tx *sql.Tx, uriHash, chocoChip string) (string, error) {
	r, err := models.Invitations(qm.Where("uri_hash = ? and choco_chip = ?", uriHash, chocoChip)).One(*ctx, tx)
	redirectUri := r.RedirectURI

	if err != nil {
		return redirectUri, err
	}
	r.AuthorisedDatetime = null.TimeFrom(time.Now())
	r.Authorised = true

	_, err = r.Update(*ctx, tx, boil.Infer())

	return redirectUri, err
}

func IsAuthorisedChocoChip(ctx *context.Context, tx *sql.Tx, chocoChip string) error {
	r, err := models.Invitations(qm.Where("choco_chip = ?", chocoChip)).One(*ctx, tx)
	if err != nil {
		/* 存在しないMLだった(Oneはレコードが存在しないと error 扱い) */
		lg.Debugf("no record found with %s", chocoChip)
		return err
	}

	if r.Authorised == false {
		return errors.New(fmt.Sprintf("%s is not authorised", chocoChip))
	}
	return nil
}

// ThrottleLimitCheck
// 同じメールアドレスに対してのマジックリンクは10分に3回まで
func ThrottleLimitCheck(ctx *context.Context, tx *sql.Tx, email string) error {
	/* mail_account にメールが登録されているかチェックしつつID取得 */
	ma, err := models.MailAccounts(qm.Where("mail_address = ?", email)).One(*ctx, tx)
	if err != nil {
		/* 未登録のメールアドレスならマジックリンクも存在するはずがないのでOK */
		return nil
	}

	tenMinAgo := time.Now().Add(-time.Minute * 10)
	r, err := models.Invitations(qm.Where("mail_account_id = ? and created_at > ?", ma.ID, tenMinAgo)).All(*ctx, tx)
	//r, err := models.Invitations(qm.Where("mail_address = ? and created_at > ?", email, tenMinAgo)).All(*ctx, tx)
	if err != nil {
		return err
	}
	if len(r) >= 100 {
		return errors.New("exceed magic-link generate throttle")
	}
	return nil
}

func GetAccountFromChocoChip(s *sessions.Session, ctx context.Context, tx *sql.Tx) (*models.Account, *models.MailAccount, error) {
	lg.Debug("authenticate.go L204")
	var (
		ma *models.MailAccount
		a  *models.Account
	)
	cc := s.Values[sess.SvNameChocochip()]
	if cc == nil {
		return a, ma, errors.New("choco_chip not found")
	}
	chocoChip := cc.(string)
	i, err := models.Invitations(qm.Where("choco_chip = ?", chocoChip)).One(ctx, tx)
	if err != nil {
		lg.Errorf("choco_chip is not found in invitation: %s", chocoChip)
		return a, ma, err
	}
	maId := i.MailAccountID
	ma, err = models.MailAccounts(qm.Where("id = ?", maId)).One(ctx, tx)
	if err != nil {
		lg.Errorf("mail_account not found: %d", maId)
		return a, ma, err
	}

	aId := ma.AccountID
	a, err = models.Accounts(qm.Where("id = ?", aId)).One(ctx, tx)
	if err != nil {
		lg.Errorf("account not found: %d", aId)
		return a, ma, err
	}

	return a, ma, nil
}

func InitInvitation(uriHash string, chocoChip string, realIP string, redirectURI string) models.Invitation {
	/* dbに登録 */
	invitation := models.Invitation{
		URIHash:            uriHash,
		ChocoChip:          chocoChip,
		IPAddress:          realIP,
		RedirectURI:        redirectURI,
		ExpiredDatetime:    time.Now().Add(3 * time.Hour),
		AuthorisedDatetime: null.Time{},
	}
	return invitation
}

func FetchMailAccountByID(ctx context.Context, tx *sql.Tx, a *models.Account) (models.MailAccountSlice, error) {
	mailAccounts, err := models.MailAccounts(qm.Where("account_id = ?", a.ID)).All(ctx, tx)
	return mailAccounts, err
}
