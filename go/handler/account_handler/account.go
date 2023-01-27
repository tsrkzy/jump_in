package account_handler

import (
	"context"
	"database/sql"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/tsrkzy/jump_in/database"
	"github.com/tsrkzy/jump_in/helper"
	"github.com/tsrkzy/jump_in/helper/validate"
	"github.com/tsrkzy/jump_in/logic/account_logic"
	"github.com/tsrkzy/jump_in/logic/authenticate_logic"
	"github.com/tsrkzy/jump_in/logic/lg"
	"github.com/tsrkzy/jump_in/sess"
	"github.com/tsrkzy/jump_in/types/account_types"
	"github.com/tsrkzy/jump_in/types/entity"
	"github.com/tsrkzy/jump_in/types/response_types"
	"net/http"
)

func UpdateName() echo.HandlerFunc {
	return func(c echo.Context) error {
		/* request */
		r := &account_types.UpdateNameRequest{}
		err := c.Bind(r)
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusBadRequest, response_types.Errors{})
		}

		accountId, err := helper.StrToID(r.AccountID)
		if err != nil {
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

		nr := &account_types.UpdateNameResponse{}
		ctx := context.Background()
		/* open DB Tx */
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			return sess.Open(c, myDB, func(s *sessions.Session) error {
				/* セッションストアからアカウントを取得 */
				a, _, err := authenticate_logic.GetAccountFromChocoChip(s, ctx, tx)
				if err != nil {
					return err
				}

				/* 認証先とリクエストのアカウントIDが異なったらNG */
				if a.ID != accountId {
					return response_types.NewErrorSeed(http.StatusUnauthorized, "認証に失敗しました")
				}

				/* アカウント名更新 */
				a.Name = r.Name
				err = account_logic.UpdateAccount(err, a, ctx, tx)
				if err != nil {
					return err
				}

				/* レスポンス作成 */
				mailAccounts, err := authenticate_logic.FetchMailAccountByID(ctx, tx, a)
				if err != nil {
					return err
				}
				maList := make([]entity.MailAccount, 0)
				for _, ma := range mailAccounts {
					maList = append(maList, *entity.CreateMailAccount(ma))
				}

				nr = &account_types.UpdateNameResponse{
					Account:      *entity.CreateAccount(a),
					MailAccounts: maList,
				}

				helper.Mask(nr)

				return nil
			})
		})
		if err != nil {
			if es, ok := err.(response_types.ErrorSeed); ok {
				return c.JSON(es.Code, es.Msg)
			}
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}

		return c.JSON(http.StatusOK, nr)
	}
}
