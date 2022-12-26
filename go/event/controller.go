package event

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/tsrkzy/jump_in/authenticate"
	"github.com/tsrkzy/jump_in/database"
	"github.com/tsrkzy/jump_in/helper"
	"github.com/tsrkzy/jump_in/lg"
	"github.com/tsrkzy/jump_in/models"
	"github.com/tsrkzy/jump_in/response"
	"github.com/tsrkzy/jump_in/sess"
	"github.com/tsrkzy/jump_in/validate"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
)

func List() echo.HandlerFunc {
	return func(c echo.Context) error {
		/* request */
		r := &ListRequest{}
		err := c.Bind(r)
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusBadRequest, response.Errors{})
		}

		accountId, err := helper.StrToID(r.AccountId)
		if err != nil {
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

		lr := ListResponse{}
		ctx := context.Background()
		/* open DB Tx */
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			return sess.Open(c, myDB, func(s *sessions.Session) error {
				/* セッションストアからアカウントを取得 */
				a, _, err := authenticate.GetAccountFromChocoChip(s, ctx, tx)
				if err != nil {
					return err
				}

				/* 認証先とリクエストのアカウントIDが異なったらNG */
				if a.ID != accountId {
					return response.NewErrorSeed(http.StatusUnauthorized, "認証に失敗しました")
				}

				aId := a.ID

				/* アカウントが作成した event 最大5件 */
				eventOwns, err := fetchEventOwns(ctx, tx, aId, 5)
				if err != nil {
					return err
				}
				lr.EventsOwns = eventOwns

				/* アカウントが join している event 最大5件 */
				eventJoins, err := fetchEventJoins(ctx, tx, aId, 5)
				if err != nil {
					return err
				}
				lr.EventsJoins = eventJoins

				/* 0 (where id in でマッチしない値) で初期化 */
				eIdExclude := make([]interface{}, 0)
				eIdExclude = append(eIdExclude, 0)
				for _, e := range eventOwns {
					eIdExclude = append(eIdExclude, e.ID)
				}
				for _, e := range eventJoins {
					eIdExclude = append(eIdExclude, e.ID)
				}
				lg.Debug(eIdExclude)

				/* 新着のイベントで、最新の最大10件 */
				eventRunning, err := fetchNewEventsWithout(ctx, tx, eIdExclude)
				if err != nil {
					return err
				}

				lr.EventsRunning = eventRunning

				return nil
			})
		})

		if err != nil {
			return c.JSON(http.StatusInternalServerError, response.Errors{})
		}
		return c.JSON(http.StatusOK, lr)
	}
}

func Detail() echo.HandlerFunc {
	return func(c echo.Context) error {
		r := &DetailRequest{}
		err := c.Bind(r)
		if err != nil {
			lg.Error(err)
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

		var dr *DetailResponse

		/* 詳細作成 */
		ctx := context.Background()
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			dr, err = getDetail(ctx, tx, r.EventId)
			return err
		})
		if err != nil {
			lg.Error(err)
			if es, ok := err.(response.ErrorSeed); ok {
				return c.JSON(es.Code, response.ErrorGen(es.Msg))
			}
			return c.JSON(http.StatusInternalServerError, response.Errors{})
		}

		return c.JSON(http.StatusOK, dr)
	}
}

func getDetail(ctx context.Context, tx *sql.Tx, eventId string) (*DetailResponse, error) {
	e, err := models.Events(qm.Where("id = ?", eventId)).One(ctx, tx)
	if err != nil {
		msg := fmt.Sprintf("イベントが見つかりません: %s", eventId)
		return nil, response.NewErrorSeed(http.StatusNotFound, msg)
	}

	owner, err := getOwner(ctx, tx, eventId)
	if err != nil {
		return nil, err
	}
	participants, err := getParticipants(ctx, tx, eventId)
	if err != nil {
		return nil, err
	}

	event := CreateEvent(e)

	dr := &DetailResponse{
		Event:        *event,
		Owner:        owner,
		Participants: participants,
	}

	return dr, nil
}

func Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		r := &CreateRequest{}
		err := c.Bind(r)
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusBadRequest, response.Errors{})
		}

		accountId, err := helper.StrToID(r.AccountId)
		if err != nil {
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

		res := &CreateResponse{}

		ctx := context.Background()
		/* open DB Tx */
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			return sess.Open(c, myDB, func(s *sessions.Session) error {
				/* セッションストアからアカウントIDを逆引き */
				a, _, err := authenticate.GetAccountFromChocoChip(s, ctx, tx)
				if err != nil {
					return err
				}

				/* 認証先とリクエストのアカウントIDが異なったらNG */
				if a.ID != accountId {
					return response.NewErrorSeed(http.StatusUnauthorized, "認証に失敗しました")
				}

				/* event_group が存在しなければ作成しておく */
				eName := r.Name
				eg, err := createEventGroup(ctx, tx, eName)
				if err != nil {
					return err
				}

				e := &models.Event{
					Name:         eName,
					EventGroupID: eg.ID,
					AccountID:    a.ID,
				}
				err = e.Insert(ctx, tx, boil.Infer())
				if err != nil {
					return err
				}

				event := CreateEvent(e)

				res = &CreateResponse{*event}

				return nil
			})
		})
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response.Errors{})
		}

		return c.JSON(http.StatusOK, res)
	}
}

func Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		r := &UpdateRequest{}
		err := c.Bind(r)
		if err != nil {
			lg.Error(err)
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

		ur := &UpdateResponse{}

		ctx := context.Background()
		/* open DB Tx */
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			return sess.Open(c, myDB, func(session *sessions.Session) error {
				eId := r.ID
				e, err := models.Events(qm.Where("id = ?", eId)).One(ctx, tx)
				if err != nil {
					return c.JSON(http.StatusNotFound, response.Errors{})
				}

				ur.Event = *CreateEvent(e)

				return nil
			})
		})

		return c.JSON(http.StatusOK, ur)
	}
}

func Attend() echo.HandlerFunc {
	return func(c echo.Context) error {
		r := &AttendRequest{}
		err := c.Bind(r)
		if err != nil {
			lg.Error(err)
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

		/* 参加予定の event を取得 */
		eId, err := helper.StrToID(r.EventId)
		if err != nil {
			return c.JSON(http.StatusNotFound, response.Errors{})
		}
		dr := &DetailResponse{}

		ctx := context.Background()
		/* open DB Tx */
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			return sess.Open(c, myDB, func(session *sessions.Session) error {

				/* セッションからアカウントIDを取得 */
				a, _, err := authenticate.GetAccountFromChocoChip(session, ctx, tx)
				if err != nil {
					return err
				}
				aId := a.ID

				_, err = models.Events(qm.Where("id = ?", eId)).One(ctx, tx)
				if err != nil {
					return response.NewErrorSeed(http.StatusNotFound, fmt.Sprintf("イベントが存在しません: %d", eId))
				}

				/* Attend から参加状況を取得 */
				exA, err := models.Attends(qm.Where("event_id = ? and account_id = ?", eId, aId)).Exists(ctx, tx)
				if err != nil {
					return err
				}

				if !exA {
					/* 未参加ならば Attend を作成 */
					att := models.Attend{
						AccountID: aId,
						EventID:   eId,
					}

					err = att.Insert(ctx, tx, boil.Infer())
					if err != nil {
						return err
					}
				}

				/* 詳細作成 */
				dr, err = getDetail(ctx, tx, r.EventId)

				return err
			})
		})
		if err != nil {
			if es, ok := err.(response.ErrorSeed); ok {
				return c.JSON(es.Code, es.Msg)
			}
			return c.JSON(http.StatusInternalServerError, response.Errors{})
		}

		ar := AttendResponse{*dr}

		return c.JSON(http.StatusOK, ar)
	}
}

func Leave() echo.HandlerFunc {
	return func(c echo.Context) error {
		r := &LeaveRequest{}
		err := c.Bind(r)
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusBadRequest, response.Errors{})
		}

		/* validation */
		if err := c.Validate(r); err != nil {
			vErr := validate.ErrorIntoJson(err)
			return c.JSON(http.StatusBadRequest, vErr)
		}
		/* 参加取消予定の event を取得 */
		eId, err := helper.StrToID(r.EventId)
		if err != nil {
			return c.JSON(http.StatusNotFound, response.Errors{})
		}

		/* db接続 */
		myDB, err := database.Open()
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response.Errors{})
		}

		dr := &DetailResponse{}

		ctx := context.Background()
		/* open DB Tx */
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			return sess.Open(c, myDB, func(session *sessions.Session) error {

				/* セッションからアカウントIDを取得 */
				a, _, err := authenticate.GetAccountFromChocoChip(session, ctx, tx)
				if err != nil {
					return err
				}
				aId := a.ID

				/* 参加取り消し予定の event を取得 */
				_, err = models.Events(qm.Where("id = ?", eId)).One(ctx, tx)
				if err != nil {
					return response.NewErrorSeed(http.StatusNotFound, fmt.Sprintf("イベントが存在しません: %d", eId))
				}

				/* Attend から参加状況を取得 */
				exA, err := models.Attends(qm.Where("event_id = ? and account_id = ?", eId, aId)).Exists(ctx, tx)
				if err != nil {
					return err
				}

				if exA {
					_, err = models.Attends(qm.Where("event_id = ? and account_id = ?", eId, aId)).DeleteAll(ctx, tx)
					if err != nil {
						return err
					}
				}

				/* 詳細作成 */
				dr, err = getDetail(ctx, tx, r.EventId)

				return nil
			})
		})
		if err != nil {
			if es, ok := err.(response.ErrorSeed); ok {
				return c.JSON(es.Code, es.Msg)
			}
			return c.JSON(http.StatusInternalServerError, response.Errors{})
		}

		lr := LeaveResponse{*dr}

		return c.JSON(http.StatusOK, lr)
	}
}
