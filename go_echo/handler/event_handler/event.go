package event_handler

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
	"github.com/tsrkzy/jump_in/logic/event_logic"
	"github.com/tsrkzy/jump_in/logic/lg"
	"github.com/tsrkzy/jump_in/sess"
	"github.com/tsrkzy/jump_in/types/entity"
	"github.com/tsrkzy/jump_in/types/event_types"
	"github.com/tsrkzy/jump_in/types/response_types"
	"net/http"
)

func List() echo.HandlerFunc {
	return func(c echo.Context) error {
		/* request */
		r := &event_types.ListRequest{}
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

		lr := event_types.ListResponse{}
		ctx := context.Background()
		/* open DB Tx */
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			return sess.Open(c, myDB, func(s *sessions.Session) error {
				/* セッションストアからアカウントを取得 */
				a, _, _, err := authenticate_logic.GetAccountFromChocoChip(s, ctx, tx)
				if err != nil {
					return err
				}

				/* 認証先とリクエストのアカウントIDが異なったらNG */
				if a.ID != accountId {
					return response_types.NewErrorSeed(http.StatusUnauthorized, "認証に失敗しました")
				}

				aId := a.ID

				/* アカウントが作成した event 最大5件 */
				eventOwns, err := event_logic.FetchEventOwns(ctx, tx, aId, 5)
				if err != nil {
					return err
				}
				lr.EventsOwns = eventOwns

				/* アカウントが join している event 最大5件 */
				eventJoins, err := event_logic.FetchEventJoins(ctx, tx, aId, 5)
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

				/* ↑のイベント以外で、公開中のイベント 新着最大10件 */
				eventRunning, err := event_logic.FetchNewEventsWithout(ctx, tx, eIdExclude)
				if err != nil {
					return err
				}

				lr.EventsRunning = eventRunning

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

		return c.JSON(http.StatusOK, lr)
	}
}

func Detail() echo.HandlerFunc {
	return func(c echo.Context) error {
		r := &event_types.DetailRequest{}
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

		/* db接続 */
		myDB, err := database.Open()
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}

		var dr *event_types.DetailResponse

		/* 詳細作成 */
		ctx := context.Background()
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			ed, err := event_logic.GetDetail(ctx, tx, r.EventID)
			dr = &event_types.DetailResponse{EventDetail: *ed}

			return err
		})
		if err != nil {
			if es, ok := err.(response_types.ErrorSeed); ok {
				return c.JSON(es.Code, response_types.ErrorGen(es.Msg))
			}
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}

		return c.JSON(http.StatusOK, dr)
	}
}

func Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		r := &event_types.CreateRequest{}
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

		res := &event_types.CreateResponse{}

		ctx := context.Background()
		/* open DB Tx */
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			return sess.Open(c, myDB, func(s *sessions.Session) error {
				/* セッションストアからアカウントIDを逆引き */
				a, _, _, err := authenticate_logic.GetAccountFromChocoChip(s, ctx, tx)
				if err != nil {
					return err
				}

				/* 認証先とリクエストのアカウントIDが異なったらNG */
				if a.ID != accountId {
					return response_types.NewErrorSeed(http.StatusUnauthorized, "認証に失敗しました")
				}

				/* event_group が存在しなければ作成しておく */
				eName := r.Name
				eg, err := event_logic.CreateEventGroup(ctx, tx, eName)
				if err != nil {
					return err
				}

				eDescription := r.Description
				e, err := event_logic.InsertEvent(ctx, tx, eName, eDescription, eg, a)
				if err != nil {
					return err
				}

				ee := entity.CreateEvent(e)

				res = &event_types.CreateResponse{*ee}

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

		return c.JSON(http.StatusOK, res)
	}
}

func UpdateName() echo.HandlerFunc {
	return func(c echo.Context) error {
		r := &event_types.UpdateNameRequest{}
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

		/* db接続 */
		myDB, err := database.Open()
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}

		ur := &event_types.UpdateNameResponse{}

		ctx := context.Background()
		/* open DB Tx */
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			return sess.Open(c, myDB, func(session *sessions.Session) error {
				eId, err := helper.StrToID(r.EventID)
				if err != nil {
					return err
				}
				e, err := event_logic.FetchEventByID(ctx, tx, eId)
				if err != nil {
					return err
				}

				e.Name = r.EventName
				err = event_logic.UpdateEvent(ctx, tx, e)
				if err != nil {
					return err
				}

				ed, err := event_logic.GetDetail(ctx, tx, r.EventID)
				ur = &event_types.UpdateNameResponse{EventDetail: *ed}

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

		return c.JSON(http.StatusOK, ur)
	}
}

func UpdateDescription() echo.HandlerFunc {
	return func(c echo.Context) error {
		r := &event_types.UpdateDescriptionRequest{}
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

		/* db接続 */
		myDB, err := database.Open()
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}

		ur := &event_types.UpdateDescriptionResponse{}

		ctx := context.Background()
		/* open DB Tx */
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			return sess.Open(c, myDB, func(session *sessions.Session) error {
				eId, err := helper.StrToID(r.EventID)
				if err != nil {
					return err
				}
				e, err := event_logic.FetchEventByID(ctx, tx, eId)
				if err != nil {
					return err
				}

				e.Description = r.Description
				err = event_logic.UpdateEvent(ctx, tx, e)
				if err != nil {
					return err
				}

				ed, err := event_logic.GetDetail(ctx, tx, r.EventID)
				ur = &event_types.UpdateDescriptionResponse{EventDetail: *ed}

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

		return c.JSON(http.StatusOK, ur)
	}
}

func UpdateOpen() echo.HandlerFunc {
	return func(c echo.Context) error {
		r := &event_types.UpdateOpenRequest{}
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

		/* db接続 */
		myDB, err := database.Open()
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}

		ur := &event_types.UpdateOpenResponse{}

		ctx := context.Background()
		/* open DB Tx */
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			return sess.Open(c, myDB, func(session *sessions.Session) error {
				eId, err := helper.StrToID(r.EventID)
				if err != nil {
					return err
				}
				e, err := event_logic.FetchEventByID(ctx, tx, eId)
				if err != nil {
					return err
				}

				e.IsOpen = r.IsOpen
				err = event_logic.UpdateEvent(ctx, tx, e)
				if err != nil {
					return err
				}

				ed, err := event_logic.GetDetail(ctx, tx, r.EventID)
				ur = &event_types.UpdateOpenResponse{EventDetail: *ed}

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

		return c.JSON(http.StatusOK, ur)
	}
}

func Attend() echo.HandlerFunc {
	return func(c echo.Context) error {
		r := &event_types.AttendRequest{}
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

		/* db接続 */
		myDB, err := database.Open()
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}

		/* 参加予定の event を取得 */
		eId, err := helper.StrToID(r.EventID)
		if err != nil {
			return c.JSON(http.StatusNotFound, response_types.Errors{})
		}
		dr := &event_types.DetailResponse{}

		ctx := context.Background()
		/* open DB Tx */
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			return sess.Open(c, myDB, func(session *sessions.Session) error {

				/* セッションからアカウントIDを取得 */
				a, _, _, err := authenticate_logic.GetAccountFromChocoChip(session, ctx, tx)
				if err != nil {
					return err
				}
				aId := a.ID

				_, err = event_logic.FetchEventByID(ctx, tx, eId)
				if err != nil {
					return response_types.NewErrorSeed(http.StatusNotFound, fmt.Sprintf("イベントが存在しません: %d", eId))
				}

				/* Attend から参加状況を取得 */
				attend, err := event_logic.FetchAttendByID(ctx, tx, aId, eId)

				comment := r.Comment

				if err == nil && attend != nil {
					/* 参加済みの場合、コメントが更新されていればUPDATE (UPSERT) */
					if attend.Comment != comment {
						lg.Debug("already exists, do UPDATE.")
						attend.Comment = comment
						err := event_logic.UpsertAttend(attend, ctx, tx)
						if err != nil {
							return err
						}
					} else {
						lg.Debug("already exists, do NOTHING.")
					}
				} else {
					/* 未参加ならば Attend を作成 */
					_, err = event_logic.CreateAttend(ctx, tx, aId, eId, comment)
					if err != nil {
						return err
					}
				}

				/* 詳細作成 */
				ed, err := event_logic.GetDetail(ctx, tx, r.EventID)
				dr = &event_types.DetailResponse{EventDetail: *ed}

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

		ar := event_types.AttendResponse{*dr}

		return c.JSON(http.StatusOK, ar)
	}
}

func Leave() echo.HandlerFunc {
	return func(c echo.Context) error {
		r := &event_types.LeaveRequest{}
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
		/* 参加取消予定の event を取得 */
		eId, err := helper.StrToID(r.EventID)
		if err != nil {
			return c.JSON(http.StatusNotFound, response_types.Errors{})
		}

		/* db接続 */
		myDB, err := database.Open()
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}

		dr := &event_types.DetailResponse{}

		ctx := context.Background()
		/* open DB Tx */
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			return sess.Open(c, myDB, func(session *sessions.Session) error {

				/* セッションからアカウントIDを取得 */
				a, _, _, err := authenticate_logic.GetAccountFromChocoChip(session, ctx, tx)
				if err != nil {
					return err
				}
				aId := a.ID

				/* 参加取り消し予定の event を取得 */
				_, err = event_logic.FetchEventByID(ctx, tx, eId)
				if err != nil {
					return response_types.NewErrorSeed(http.StatusNotFound, fmt.Sprintf("イベントが存在しません: %d", eId))
				}

				/* Attend から参加状況を取得 */
				exA, err := event_logic.AttendExists(ctx, tx, aId, eId)
				if err != nil {
					return err
				}

				if exA {
					err = event_logic.DeleteAttendByIDs(ctx, tx, aId, eId)
					if err != nil {
						return err
					}
				}

				/* 詳細作成 */
				ed, err := event_logic.GetDetail(ctx, tx, r.EventID)
				dr = &event_types.DetailResponse{EventDetail: *ed}

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

		lr := event_types.LeaveResponse{*dr}

		return c.JSON(http.StatusOK, lr)
	}
}

func Certify() echo.HandlerFunc {
	return func(c echo.Context) error {
		r := &event_types.CertifyRequest{}
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

		/* db接続 */
		myDB, err := database.Open()
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}

		dr := &event_types.DetailResponse{}

		ctx := context.Background()
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			return sess.Open(c, myDB, func(session *sessions.Session) error {

				/* セッションからアカウントIDを取得 */
				a, _, _, err := authenticate_logic.GetAccountFromChocoChip(session, ctx, tx)
				if err != nil {
					return err
				}

				/* 認証先とリクエストのアカウントIDが異なったらNG */
				if a.ID != accountId {
					return response_types.NewErrorSeed(http.StatusUnauthorized, "認証に失敗しました")
				}

				_, err = event_logic.CertifyEvent(ctx, tx, eId, r.Certify)
				if err != nil {
					return err
				}

				/* 詳細作成 */
				ed, err := event_logic.GetDetail(ctx, tx, r.EventID)
				dr = &event_types.DetailResponse{EventDetail: *ed}

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

		cr := &event_types.CertifyResponse{*dr}

		return c.JSON(http.StatusOK, cr)
	}
}
