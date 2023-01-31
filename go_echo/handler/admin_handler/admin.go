package admin_handler

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/tsrkzy/jump_in/database"
	"github.com/tsrkzy/jump_in/helper"
	"github.com/tsrkzy/jump_in/helper/validate"
	"github.com/tsrkzy/jump_in/logic/authenticate_logic"
	"github.com/tsrkzy/jump_in/logic/lg"
	"github.com/tsrkzy/jump_in/models"
	"github.com/tsrkzy/jump_in/sess"
	"github.com/tsrkzy/jump_in/types/admin_types"
	"github.com/tsrkzy/jump_in/types/response_types"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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
				invitation, err := models.Invitations(qm.Where("choco_chip = ?", chocochip)).One(ctx, tx)
				if err != nil || invitation == nil {
					return response_types.ErrorSeed{Code: http.StatusUnauthorized}
				}

				/* adminが既に存在するなら終了 */
				admin, err := models.Administrators(qm.Where("account_id = ? and invitation_id = ?", accountId, invitation.ID)).One(ctx, tx)
				if err == nil && admin != nil {
					return nil
				}
				/* pass_hashの突合 */
				pass_hash := r.PassHash
				pass := "fushianasan"
				true_hash := fmt.Sprintf("%x", sha256.Sum256([]byte(pass)))
				if pass_hash != true_hash {
					lg.Debug("wrong password")
					return response_types.ErrorSeed{Code: http.StatusBadRequest, Msg: "パスワードが異なります"}
				}

				/* adminにInsert */
				admin = &models.Administrator{
					AccountID:    accountId,
					InvitationID: invitation.ID,
				}
				return admin.Insert(ctx, tx, boil.Infer())
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
