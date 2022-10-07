package event

import (
	"context"
	"database/sql"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/tsrkzy/jump_in/authenticate"
	"github.com/tsrkzy/jump_in/database"
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
				aId := a.ID

				/* アカウントが作成した event 最大5件 */
				eOwns, err := models.Events(
					qm.Where("account_id = ?", aId),
					qm.Limit(5),
					qm.OrderBy("created_at desc"),
				).All(ctx, tx)
				if err != nil {
					return err
				}

				eventOwns := make([]Event, 0)
				for _, e := range eOwns {
					ee := Event{*e}
					eventOwns = append(eventOwns, ee)
				}

				lr.EventsOwns = eventOwns

				/* アカウントが join している event 最大5件 */
				// attend
				at, err := models.Attends(
					qm.Where("account_id = ?", aId),
					qm.Limit(5),
					qm.OrderBy("created_at desc"),
				).All(ctx, tx)
				eIdList := make([]interface{}, 0)
				for _, a := range at {
					eId := a.EventID
					eIdList = append(eIdList, int(eId))
				}

				eJoins, err := models.Events(
					qm.WhereIn("id in ?", eIdList...),
					qm.Limit(3),
				).All(ctx, tx)
				if err != nil {
					return err
				}

				eventJoins := make([]Event, 0)
				for _, e := range eJoins {
					ee := Event{*e}
					eventJoins = append(eventJoins, ee)
				}

				lr.EventsJoins = eventJoins

				/* 0 (where id in でマッチしない値) で初期化 */
				eIdExclude := make([]interface{}, 0)
				eIdExclude = append(eIdExclude, 0)
				for _, e := range eOwns {
					eIdExclude = append(eIdExclude, e.ID)
				}
				for _, e := range eJoins {
					eIdExclude = append(eIdExclude, e.ID)
				}
				lg.Debug(eIdExclude)

				/* 新着のイベントで、最新の最大10件 */
				eList, err := models.Events(
					qm.WhereNotIn("id not in ?", eIdExclude...),
					qm.Limit(10),
					qm.OrderBy("created_at DESC"),
				).All(ctx, tx)
				if err != nil {
					return err
				}

				eventList := make([]Event, 0)
				for _, e := range eList {
					ee := Event{*e}
					eventList = append(eventList, ee)
				}

				lr.EventsRunning = eventList

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
		lg.Info("r")
		lg.Info(r)
		lg.Info(r.EventId)

		/* db接続 */
		myDB, err := database.Open()
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response.Errors{})
		}

		notFound := false
		dr := &DetailResponse{}

		ctx := context.Background()
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			e, err := models.Events(qm.Where("id = ?", r.EventId)).One(ctx, tx)
			if err != nil {
				notFound = true
				return err
			}
			dr = &DetailResponse{Event{*e}}

			return nil
		})
		if err != nil {
			lg.Error(err)
			if notFound {
				return c.JSON(http.StatusNotFound, response.Errors{})
			}
			return c.JSON(http.StatusInternalServerError, response.Errors{})
		}

		return c.JSON(http.StatusOK, dr)
	}
}

func Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		r := &CreateRequest{}
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

		ctx := context.Background()
		/* open DB Tx */
		err = myDB.Tx(ctx, func(tx *sql.Tx) error {
			return sess.Open(c, myDB, func(s *sessions.Session) error {

				/* event_group が存在しなければ作成しておく */
				eName := r.Name
				eg, err := createEventGroup(ctx, tx, eName)
				if err != nil {
					return err
				}

				/* セッションストアからアカウントIDを逆引き */
				a, _, err := authenticate.GetAccountFromChocoChip(s, ctx, tx)
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

				return nil
			})
		})
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response.Errors{})
		}

		res := &CreateResponse{}

		return c.JSON(http.StatusOK, res)
	}
}
