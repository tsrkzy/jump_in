package admin_handler

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/tsrkzy/jump_in/database"
	"github.com/tsrkzy/jump_in/helper"
	"github.com/tsrkzy/jump_in/helper/validate"
	"github.com/tsrkzy/jump_in/logic/administrator_logic"
	"github.com/tsrkzy/jump_in/logic/authenticate_logic"
	"github.com/tsrkzy/jump_in/logic/lg"
	"github.com/tsrkzy/jump_in/sess"
	"github.com/tsrkzy/jump_in/types/admin_types"
	"github.com/tsrkzy/jump_in/types/response_types"
	"net/http"
)

// Login
// パスワードの入力間違い = 400 BadRequest
// 他は基本的に 401 Unauthorised
func Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		/* request body */
		r := admin_types.LoginRequest{}
		if err := c.Bind(&r); err != nil {
			log.Debug(err)
			return c.JSON(http.StatusBadRequest, response_types.Errors{})
		}

		/* validation */
		if err := c.Validate(r); err != nil {
			vErr := validate.ErrorIntoJson(err)
			return c.JSON(http.StatusBadRequest, vErr)
		}

		accountId, err := helper.StrToID(r.AccountID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}

		/* db接続 */
		myDB, err := database.Open()
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}

		ctx := context.Background()
		/* Open Db Tx */
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			/* Open OpenSession Store */
			return sess.Open(c, myDB, func(s *sessions.Session) error {
				/* accountID */
				a, _, _, err := authenticate_logic.GetAccountFromChocoChip(s, ctx, tx)
				if err != nil {
					return err
				}

				/* セッションとBODYのaccountIDが等しいか検証 */
				if accountId != a.ID {
					return response_types.ErrorSeed{Code: http.StatusUnauthorized}
				}

				/* chocochipでInvitationを取得 */
				chocochip := s.Values[sess.SvNameChocochip()]
				invitation, err := authenticate_logic.FetchInvitationByChocochip(ctx, tx, chocochip.(string))
				if err != nil || invitation == nil {
					return response_types.ErrorSeed{Code: http.StatusUnauthorized}
				}

				/* adminが既に存在するなら終了 */
				admin, err := administrator_logic.FetchAdministratorByUK(ctx, tx, accountId, invitation.ID)
				if err == nil && admin != nil {
					return nil
				}

				/* pass_hashの突合 */
				pass_hash := r.PassHash
				pass := "fushianasan"
				true_hash := helper.Sha256Digest(pass)
				if pass_hash != true_hash {
					lg.Debug("wrong password")
					return response_types.ErrorSeed{Code: http.StatusBadRequest, Msg: "パスワードが異なります"}
				}

				/* adminにInsert */
				return administrator_logic.InsertAdministrator(ctx, tx, accountId, invitation.ID)
			})
		})

		if err != nil {
			if es, ok := err.(response_types.ErrorSeed); ok {
				return c.JSON(es.Code, es.Msg)
			}
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}

		lr := &admin_types.LoginResponse{
			OK: response_types.Ok(),
		}

		return c.JSON(http.StatusOK, lr)
	}
}

func Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		/* request body */
		r := admin_types.LogoutRequest{}
		if err := c.Bind(&r); err != nil {
			log.Debug(err)
			return c.JSON(http.StatusBadRequest, response_types.Errors{})
		}

		/* validation */
		if err := c.Validate(r); err != nil {
			vErr := validate.ErrorIntoJson(err)
			return c.JSON(http.StatusBadRequest, vErr)
		}

		accountId, err := helper.StrToID(r.AccountID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}

		/* db接続 */
		myDB, err := database.Open()
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}

		ctx := context.Background()
		/* Open Db Tx */
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			/* Open OpenSession Store */
			return sess.Open(c, myDB, func(s *sessions.Session) error {

				/* accountID */
				a, _, admin, err := authenticate_logic.GetAccountFromChocoChip(s, ctx, tx)
				if err != nil {
					return err
				}

				/* セッションとBODYのaccountIDが等しいか検証 */
				if accountId != a.ID {
					return response_types.ErrorSeed{Code: http.StatusUnauthorized}
				}

				if admin == nil {
					/* 管理者Tにレコードが存在しないなら終了 */
					lg.Debug("no records in administorator.")
					return nil
				}

				_, err = admin.Delete(ctx, tx)

				return err
			})
		})

		if err != nil {
			if es, ok := err.(response_types.ErrorSeed); ok {
				return c.JSON(es.Code, es.Msg)
			}
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}

		lr := &admin_types.LogoutResponse{
			OK: response_types.Ok(),
		}

		return c.JSON(http.StatusOK, lr)
	}
}

// Ad
// API(HandlerFunc)の呼び出しに管理者権限の認証をかける
// Auと一緒
func Ad(f echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := isAdmin(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, response_types.ErrorGen("管理者ログインが必要です"))
		}

		return f(c)
	}
}

func isAdmin(c echo.Context) error {

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
			/* chocochipでInvitationとAdminが取得できるか */
			cc := s.Values[sess.SvNameChocochip()]
			if cc == nil {
				return errors.New("choco_chip not found")
			}
			chocoChip := cc.(string)
			lg.Debugf("choco_chip: %s", chocoChip)

			_, _, admin, err := authenticate_logic.GetAccountFromChocoChip(s, ctx, tx)
			if err != nil {
				return err
			}

			if admin == nil {
				return errors.New("administrator not found")
			}

			return nil
		})
		return err
	})

	if err != nil {
		lg.Debug(err)
		return err
	}
	return nil
}
