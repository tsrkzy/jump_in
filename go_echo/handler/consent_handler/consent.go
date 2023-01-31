package consent_handler

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/tsrkzy/jump_in/database"
	"github.com/tsrkzy/jump_in/helper"
	"github.com/tsrkzy/jump_in/helper/validate"
	"github.com/tsrkzy/jump_in/logic/authenticate_logic"
	"github.com/tsrkzy/jump_in/logic/consent_logic"
	"github.com/tsrkzy/jump_in/logic/event_logic"
	"github.com/tsrkzy/jump_in/logic/lg"
	"github.com/tsrkzy/jump_in/sess"
	"github.com/tsrkzy/jump_in/types/consent_types"
	"github.com/tsrkzy/jump_in/types/response_types"
	"net/http"
)

func Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		/* request */
		r := &consent_types.CreateRequest{}
		err := c.Bind(r)
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusBadRequest, response_types.Errors{})
		}

		/* validation */
		if err := c.Validate(r); err != nil {
			vErr := validate.ErrorIntoJson(err)
			return c.JSON(http.StatusBadRequest, vErr)
		}

		eId, err := helper.StrToID(r.EventID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response_types.Errors{})
		}

		accountId, err := helper.StrToID(r.AccountID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response_types.Errors{})
		}

		message := r.Message

		/* db接続 */
		myDB, err := database.Open()
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}

		cr := &consent_types.CreateResponse{}

		ctx := context.Background()
		/* open DB Tx */
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			return sess.Open(c, myDB, func(s *sessions.Session) error {
				/* セッションからアカウントIDを取得 */
				a, _, admin, err := authenticate_logic.GetAccountFromChocoChip(s, ctx, tx)
				if err != nil {
					return err
				}

				/* 認証先とリクエストのアカウントIDが異なったらNG */
				if a.ID != accountId {
					return response_types.NewErrorSeed(http.StatusUnauthorized, "認証に失敗しました")
				}

				/* イベント存在チェック */
				e, err := event_logic.FetchEventByID(ctx, tx, eId)
				if err != nil {
					return response_types.NewErrorSeed(http.StatusNotFound, fmt.Sprintf("イベントが存在しません: %d", eId))
				}

				/* イベントオーナー */
				ownerId := e.AccountID

				err = consent_logic.InsertConsent(ctx, tx, eId, admin.AccountID, ownerId, message)
				if err != nil {
					return err
				}

				ed, err := event_logic.GetDetail(ctx, tx, r.EventID)
				cr = &consent_types.CreateResponse{EventDetail: *ed}

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

		return c.JSON(http.StatusOK, cr)
	}
}

func Accept() echo.HandlerFunc {
	return func(c echo.Context) error {
		/* request */
		r := &consent_types.AcceptRequest{}
		err := c.Bind(r)
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusBadRequest, response_types.Errors{})
		}

		/* validation */
		if err := c.Validate(r); err != nil {
			vErr := validate.ErrorIntoJson(err)
			return c.JSON(http.StatusBadRequest, vErr)
		}

		eId, err := helper.StrToID(r.EventID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response_types.Errors{})
		}

		accountId, err := helper.StrToID(r.AccountID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response_types.Errors{})
		}

		consentId, err := helper.StrToID(r.ConsentID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response_types.Errors{})
		}

		ar := &consent_types.AcceptResponse{}

		/* db接続 */
		myDB, err := database.Open()
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}

		ctx := context.Background()
		/* open DB Tx */
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			return sess.Open(c, myDB, func(s *sessions.Session) error {
				/* セッションからアカウントIDを取得 */
				a, _, _, err := authenticate_logic.GetAccountFromChocoChip(s, ctx, tx)
				if err != nil {
					return err
				}

				/* 認証先とリクエストのアカウントIDが異なったらNG */
				if a.ID != accountId {
					return response_types.NewErrorSeed(http.StatusUnauthorized, "認証に失敗しました")
				}

				/* イベント存在チェック */
				e, err := event_logic.FetchEventByID(ctx, tx, eId)
				if err != nil {
					return response_types.NewErrorSeed(http.StatusNotFound, fmt.Sprintf("イベントが存在しません: %d", eId))
				}

				if accountId != e.AccountID {
					return response_types.NewErrorSeed(http.StatusBadRequest, fmt.Sprint("イベントの所有者ではありません"))
				}

				/* 同意書 */
				consent, err := consent_logic.FetchConsentByUK(ctx, tx, consentId, accountId, e.ID)
				if err != nil {
					lg.Debugf("no consent found. id: %d, %d, %d", consentId, accountId, e.ID)
					return err
				}

				err = consent_logic.UpdateConsentAccepted(ctx, tx, consent)
				if err != nil {
					return err
				}

				ed, err := event_logic.GetDetail(ctx, tx, r.EventID)
				ar = &consent_types.AcceptResponse{EventDetail: *ed}

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

		return c.JSON(http.StatusOK, ar)
	}
}
