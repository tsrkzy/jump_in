package authenticate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/tsrkzy/jump_in/database"
	"github.com/tsrkzy/jump_in/lg"
	"github.com/tsrkzy/jump_in/mail"
	"github.com/tsrkzy/jump_in/models"
	"github.com/tsrkzy/jump_in/response"
	"github.com/tsrkzy/jump_in/sess"
	"github.com/tsrkzy/jump_in/validate"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"time"
)

func Authenticate() echo.HandlerFunc {
	return func(c echo.Context) error {
		/* request body */
		r := Request{}
		if err := c.Bind(&r); err != nil {
			log.Debug(err)
			return c.JSON(http.StatusBadRequest, response.Errors{})
		}

		/* validation */
		if err := c.Validate(r); err != nil {
			vErr := validate.ErrorIntoJson(err)
			return c.JSON(http.StatusBadRequest, vErr)
		}

		/* db接続 */
		myDB, err := database.Open()
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response.Errors{})
		}

		ar := Result{}
		ctx := context.Background()
		/* open DB Tx */
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			/* Open SessionStore */
			return sess.Open(c, myDB, func(s *sessions.Session) error {
				s.Options = &sessions.Options{
					MaxAge:   60 * 5,
					HttpOnly: true,
					Secure:   true,
				}

				/* スロットル: 同じメールアドレスについて、10分に3回まで */
				if err := ThrottleLimitCheck(&ctx, tx, r.MailAddress); err != nil {
					return response.NewErrorSeed(http.StatusTooManyRequests, "時間を置いて再度リクエストしてください")
				}

				/* url_hash生成 */
				uriHash, err := GenerateHash(&ctx, tx, "uri_hash")
				if err != nil {
					return err
				}

				/* choco_chip生成 */
				chocoChip, err := GenerateHash(&ctx, tx, "choco_chip")
				if err != nil {
					return err
				}

				/* dbに登録 */
				invitation := models.Invitation{
					URIHash:            uriHash,
					ChocoChip:          chocoChip,
					IPAddress:          c.RealIP(),
					RedirectURI:        r.RedirectURI,
					ExpiredDatetime:    time.Now().Add(3 * time.Hour),
					AuthorisedDatetime: null.Time{},
				}
				if err := CreateInvitation(&ctx, tx, &invitation, r.MailAddress); err != nil {
					return err
				}

				/* メール送信 */
				ml := "http://localhost:80/api/ml/" + uriHash
				m := mail.Content{
					MailTo:  r.MailAddress,
					NameTo:  "宛先@TODO",
					Subject: "Invitation Link",
					Body:    ml,
				}
				if err := mail.SendMailSSL(&m); err != nil {
					return err
				}

				/* response */
				ar.MailAddress = r.MailAddress
				ar.URIHash = uriHash
				ar.ChocoChip = chocoChip
				ar.MagicLink = ml
				ar.IpAddress = c.RealIP()

				/* save session-cookie */
				s.Values[sess.SvNameChocochip()] = chocoChip

				return nil
			})
		})
		if err != nil {
			lg.Error(err)
			if es, ok := err.(response.ErrorSeed); ok {
				return c.JSON(es.Code, response.ErrorGen(es.Msg))
			}
			return c.JSON(http.StatusInternalServerError, response.Errors{})
		}

		return c.JSON(http.StatusOK, ar)
	}
}

func MagicLink() echo.HandlerFunc {
	return func(c echo.Context) error {
		/* クロールで露出したくないので、基本5xx系ではなく404を返す */
		uriHash := c.Param("uri_hash")

		/* db接続 */
		myDB, err := database.Open()
		if err != nil {
			lg.Error(err)
			return c.NoContent(http.StatusNotFound)
		}

		var redirectUri string

		ctx := context.Background()

		/* open DB Tx */
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			/* Open session store */
			return sess.Open(c, myDB, func(s *sessions.Session) error {

				/* chocochip の検証 */
				cc := s.Values[sess.SvNameChocochip()]
				if cc == nil {
					return errors.New("not found: session cookie")
				}
				chocoChip := cc.(string)
				lg.Debug("choco-chip ok")

				exist, err := CheckExistence(&ctx, tx, uriHash, chocoChip)
				if err != nil {
					return err
				}

				if !exist {
					return errors.New(fmt.Sprintf("hash %s does not exist", uriHash))
				}
				lg.Debug("magic-link available")

				/* 使用したurlを非活性化 */
				redirectUri, err = AuthoriseMagicLink(&ctx, tx, uriHash, chocoChip)
				if err != nil {
					return err
				}
				lg.Debug("magic-link authorised")

				/* セッション変数をログイン済みに昇格、寿命を3時間に */
				s.Options = &sessions.Options{
					MaxAge:   60 * 60 * 3,
					HttpOnly: true,
					Secure:   true,
				}
				lg.Debugf("session choco_chip promoted: %s", chocoChip)

				return nil
			})
		})
		if err != nil {
			lg.Error(err)
			return c.NoContent(http.StatusNotFound)
		}

		/* redirect with 302: Found */
		return c.Redirect(http.StatusFound, redirectUri)
	}
}

func Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		myDB, err := database.Open()
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response.Errors{})
		}
		err = sess.Open(c, myDB, func(s *sessions.Session) error {
			s.Options = &sessions.Options{
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   true,
			}
			return nil
		})
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response.Errors{})
		}
		return c.JSON(http.StatusOK, response.Ok())
	}
}

func Status() echo.HandlerFunc {
	return func(c echo.Context) error {
		/* db接続 */
		myDB, err := database.Open()
		if err != nil {
			lg.Error(err)
			return c.NoContent(http.StatusNotFound)
		}

		ctx := context.Background()
		/* Open Db Tx */
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			/* Open OpenSession Store */
			err = sess.Open(c, myDB, func(s *sessions.Session) error {
				cc := s.Values[sess.SvNameChocochip()]
				if cc == nil {
					return errors.New("choco_chip not found")
				}
				chocoChip := cc.(string)
				lg.Debugf("choco_chip: %s", chocoChip)

				return IsAuthorisedChocoChip(&ctx, tx, chocoChip)
			})
			return err
		})

		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusUnauthorized, response.Errors{})
		} else {
			return c.JSON(http.StatusOK, response.Ok())
		}
	}
}

func WhoAmI() echo.HandlerFunc {
	return func(c echo.Context) error {

		/* DB接続*/
		myDB, err := database.Open()
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response.Errors{})
		}

		wr := &WhoAmIResponse{}

		ctx := context.Background()
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			return sess.Open(c, myDB, func(session *sessions.Session) error {
				a, _, err := GetAccountFromChocoChip(session, ctx, tx)
				if err != nil {
					return err
				}

				mailAccounts, err := models.MailAccounts(qm.Where("account_id = ?", a.ID)).All(ctx, tx)
				if err != nil {
					return err
				}
				maList := make([]models.MailAccount, 0)
				for _, ma := range mailAccounts {
					maList = append(maList, *ma)
				}

				wr = &WhoAmIResponse{
					Account:      *a,
					MailAccounts: maList,
				}

				return nil
			})
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response.Errors{})
		}
		return c.JSON(http.StatusOK, wr)
	}
}
