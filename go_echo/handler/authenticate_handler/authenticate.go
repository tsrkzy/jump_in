package authenticate_handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/tsrkzy/jump_in/database"
	"github.com/tsrkzy/jump_in/helper"
	"github.com/tsrkzy/jump_in/helper/validate"
	"github.com/tsrkzy/jump_in/logic/authenticate_logic"
	"github.com/tsrkzy/jump_in/logic/lg"
	"github.com/tsrkzy/jump_in/logic/mail"
	"github.com/tsrkzy/jump_in/sess"
	"github.com/tsrkzy/jump_in/types/authenticate_types"
	"github.com/tsrkzy/jump_in/types/entity"
	"github.com/tsrkzy/jump_in/types/response_types"
	"net/http"
)

func Authenticate() echo.HandlerFunc {
	return func(c echo.Context) error {
		/* request body */
		r := authenticate_types.Request{}
		if err := c.Bind(&r); err != nil {
			log.Debug(err)
			return c.JSON(http.StatusBadRequest, response_types.Errors{})
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
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}

		m := mail.Content{}

		ar := authenticate_types.Result{}
		ctx := context.Background()
		/* open DB Tx */
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			/* Open SessionStore */
			return sess.Open(c, myDB, func(s *sessions.Session) error {
				s.Options = &sessions.Options{
					MaxAge:   60 * 5,
					HttpOnly: true,
					Secure:   true,
					Path:     "/",
				}

				/* スロットル: 同じメールアドレスについて、10分に3回まで */
				mailAddress := r.MailAddress
				if err := authenticate_logic.ThrottleLimitCheck(&ctx, tx, mailAddress); err != nil {
					return response_types.NewErrorSeed(http.StatusTooManyRequests, "時間を置いて再度リクエストしてください")
				}

				/* url_hash生成 */
				uriHash, err := authenticate_logic.GenerateHash(&ctx, tx, "uri_hash")
				if err != nil {
					return err
				}

				/* choco_chip生成 */
				chocoChip, err := authenticate_logic.GenerateHash(&ctx, tx, "choco_chip")
				if err != nil {
					return err
				}

				invitation := authenticate_logic.InitInvitation(uriHash, chocoChip, c.RealIP(), r.RedirectURI)
				if err := authenticate_logic.CreateInvitation(&ctx, tx, &invitation, mailAddress); err != nil {
					return err
				}

				/* メール送信 */
				redirectUri := r.RedirectURI

				redirectOrigin, err := helper.ExtractOriginFromURI(redirectUri)
				if err != nil {
					return err
				}
				ml := redirectOrigin + "/api/ml/" + uriHash
				m = mail.Content{
					MailTo:  mailAddress,
					NameTo:  "JumpIn参加者様",
					Subject: "JumpIn 認証用マジックリンク",
					Body:    ml,
				}

				/* response */
				ar.MailAddress = mailAddress
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
			if es, ok := err.(response_types.ErrorSeed); ok {
				return c.JSON(es.Code, response_types.ErrorGen(es.Msg))
			}
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}

		if err = mail.SendMailSSL(&m); err != nil {
			return c.JSON(http.StatusInternalServerError, response_types.ErrorGen("メール送信に失敗しました"))
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

				exist, err := authenticate_logic.CheckExistence(&ctx, tx, uriHash, chocoChip)
				if err != nil {
					return err
				}

				if !exist {
					return errors.New(fmt.Sprintf("hash %s does not exist", uriHash))
				}
				lg.Debug("magic-link available")

				/* 使用したurlを非活性化 */
				redirectUri, err = authenticate_logic.AuthoriseMagicLink(&ctx, tx, uriHash, chocoChip)
				if err != nil {
					return err
				}
				lg.Debug("magic-link authorised")

				/* セッション変数をログイン済みに昇格、寿命を3時間に */
				s.Options = &sessions.Options{
					MaxAge:   60 * 60 * 3,
					HttpOnly: true,
					Secure:   true,
					Path:     "/",
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
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}
		err = sess.Open(c, myDB, func(s *sessions.Session) error {
			s.Options = &sessions.Options{
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   true,
				Path:     "/",
			}
			return nil
		})
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}
		return c.JSON(http.StatusOK, response_types.Ok())
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

				return authenticate_logic.IsAuthorisedChocoChip(&ctx, tx, chocoChip)
			})
			return err
		})

		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusUnauthorized, response_types.Errors{})
		} else {
			return c.JSON(http.StatusOK, response_types.Ok())
		}
	}
}

func WhoAmI() echo.HandlerFunc {
	return func(c echo.Context) error {

		/* DB接続*/
		myDB, err := database.Open()
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}

		wr := &authenticate_types.WhoAmIResponse{}

		ctx := context.Background()
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			return sess.Open(c, myDB, func(s *sessions.Session) error {
				a, _, _admin, err := authenticate_logic.GetAccountFromChocoChip(s, ctx, tx)
				if err != nil {
					return err
				}

				mailAccounts, err := authenticate_logic.FetchMailAccountByID(ctx, tx, a)
				if err != nil {
					return err
				}
				maList := make([]entity.MailAccount, 0)
				for _, ma := range mailAccounts {
					maList = append(maList, *entity.CreateMailAccount(ma))
				}

				admin := &entity.Administrator{}
				if _admin != nil {
					lg.Debug("got admin")
					admin = entity.CreateAdministrator(_admin)
				} else {
					lg.Debug("got no admin")
				}

				wr = &authenticate_types.WhoAmIResponse{
					Account:       *entity.CreateAccount(a),
					Administrator: *admin,
					MailAccounts:  maList,
				}
				helper.Mask(wr)

				return nil
			})
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}
		return c.JSON(http.StatusOK, wr)
	}
}

// Au
// API(HandlerFunc)の呼び出しにログイン認証をかける
// ログイン済みならそのまま Handler を処理
// 未ログインなら 401 Unauthorized で応答
func Au(f echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := authorized(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, response_types.ErrorGen("ログインが必要です"))
		}

		return f(c)
	}
}

func authorized(c echo.Context) error {
	/* db接続 */
	myDB, err := database.Open()
	if err != nil {
		lg.Fatal(err)
		return err
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

			return authenticate_logic.IsAuthorisedChocoChip(&ctx, tx, chocoChip)
		})
		return err
	})

	if err != nil {
		lg.Debug(err)
		return err
	} else {
		return nil
	}
}
