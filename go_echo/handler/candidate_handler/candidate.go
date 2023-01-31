package candidate_handler

import (
	"context"
	"database/sql"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/tsrkzy/jump_in/database"
	"github.com/tsrkzy/jump_in/helper"
	"github.com/tsrkzy/jump_in/helper/validate"
	"github.com/tsrkzy/jump_in/logic/authenticate_logic"
	"github.com/tsrkzy/jump_in/logic/candidate_logic"
	"github.com/tsrkzy/jump_in/logic/event_logic"
	"github.com/tsrkzy/jump_in/logic/lg"
	"github.com/tsrkzy/jump_in/sess"
	"github.com/tsrkzy/jump_in/types/candidate_types"
	"github.com/tsrkzy/jump_in/types/response_types"
	"net/http"
)

func Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		/* request body */
		r := candidate_types.CreateRequest{}
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
			return c.JSON(http.StatusBadRequest, response_types.Errors{})
		}

		eId, err := helper.StrToID(r.EventID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response_types.Errors{})
		}

		/* db接続 */
		myDB, err := database.Open()
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}

		cr := candidate_types.CreateResponse{}
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

				/* イベント作成者のアカウントIDとリクエストのアカウントIDが異なったらNG */
				e, err := event_logic.FetchEventByID(ctx, tx, eId)
				if err != nil {
					return err
				} else if accountId != e.AccountID {
					return response_types.NewErrorSeed(http.StatusUnauthorized, "候補日の操作はイベントの作成者のみ可能です")
				}

				candidate, err := candidate_logic.FetchCandidateByDate(ctx, tx, eId, r.OpenAt)
				if err != nil || candidate == nil {
					_, err := candidate_logic.CreateCandidate(ctx, tx, eId, r.OpenAt)
					if err != nil {
						return err
					}
				} else {
					/* 既に存在していたら何もしない */
					lg.Debugf("candidate with eventID: %d, openAt: %s already exists.", eId, r.OpenAt)
				}

				ed, err := event_logic.GetDetail(ctx, tx, r.EventID)
				if err != nil {
					return err
				}

				cr = candidate_types.CreateResponse{EventDetail: *ed}

				return nil
			})
		})
		if err != nil {
			if es, ok := err.(response_types.ErrorSeed); ok {
				return c.JSON(es.Code, es.Msg)
			}
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}
		return c.JSON(http.StatusOK, cr)
	}
}

func Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		/* request body */
		r := candidate_types.DeleteRequest{}
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
			return c.JSON(http.StatusBadRequest, response_types.Errors{})
		}

		eId, err := helper.StrToID(r.EventID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response_types.Errors{})
		}

		cId, err := helper.StrToID(r.CandidateID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response_types.Errors{})
		}

		/* db接続 */
		myDB, err := database.Open()
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}

		dr := candidate_types.DeleteResponse{}

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

				/* イベント作成者のアカウントIDとリクエストのアカウントIDが異なったらNG */
				e, err := event_logic.FetchEventByID(ctx, tx, eId)
				if err != nil {
					return err
				} else if accountId != e.AccountID {
					return response_types.NewErrorSeed(http.StatusUnauthorized, "候補日の操作はイベントの作成者のみ可能です")
				}

				candidate, err := candidate_logic.FetchCandidateByID(ctx, tx, cId, eId)
				if err != nil || candidate == nil {
					/* 存在しなければ何もしない */
					lg.Debugf("candidate with eventID: %d, candidateID: %d does not exists.", eId, cId)
				} else {
					_, err := candidate.Delete(ctx, tx)
					if err != nil {
						return err
					}
				}

				ed, err := event_logic.GetDetail(ctx, tx, r.EventID)
				if err != nil {
					return err
				}
				dr = candidate_types.DeleteResponse{EventDetail: *ed}

				return nil
			})
		})

		if err != nil {
			if es, ok := err.(response_types.ErrorSeed); ok {
				return c.JSON(es.Code, es.Msg)
			}
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}

		return c.JSON(http.StatusOK, dr)
	}
}

func Upvote() echo.HandlerFunc {
	return func(c echo.Context) error {
		/* request body */
		r := candidate_types.UpvoteRequest{}
		if err := c.Bind(&r); err != nil {
			return c.JSON(http.StatusBadRequest, response_types.Errors{})
		}

		/* validation */
		if err := c.Validate(r); err != nil {
			vErr := validate.ErrorIntoJson(err)
			return c.JSON(http.StatusBadRequest, vErr)
		}

		lg.Debug("candidate.go L221")

		accountId, err := helper.StrToID(r.AccountID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response_types.Errors{})
		}

		eId, err := helper.StrToID(r.EventID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response_types.Errors{})
		}

		cId, err := helper.StrToID(r.CandidateID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response_types.Errors{})
		}

		lg.Debug("candidate.go L236")

		/* db接続 */
		myDB, err := database.Open()
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}

		cr := candidate_types.UpvoteResponse{}
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

				/* 候補日の存在チェック */
				exists, err := candidate_logic.CandidateExistsByID(ctx, tx, cId, eId)
				if err != nil {
					return err
				}
				if !exists {
					return response_types.NewErrorSeed(http.StatusNotFound, "候補日が存在しません")
				}

				/* 投票の存在チェック */
				exists, err = candidate_logic.VoteExistsByUK(ctx, tx, accountId, cId)
				if err != nil {
					return err
				}
				if exists {
					/* 既に存在するなら何もしない */
					lg.Debugf("vote already exists. candidate: %d, account: %d", cId, accountId)
				} else {
					_, err := candidate_logic.CreateVote(ctx, tx, accountId, cId)
					if err != nil {
						return err
					}
				}

				ed, err := event_logic.GetDetail(ctx, tx, r.EventID)

				cr = candidate_types.UpvoteResponse{EventDetail: *ed}

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
		return c.JSON(http.StatusOK, cr)
	}
}

func Downvote() echo.HandlerFunc {
	return func(c echo.Context) error {
		/* request body */
		r := candidate_types.DownvoteRequest{}
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
			return c.JSON(http.StatusBadRequest, response_types.Errors{})
		}

		eId, err := helper.StrToID(r.EventID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response_types.Errors{})
		}

		cId, err := helper.StrToID(r.CandidateID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response_types.Errors{})
		}

		/* db接続 */
		myDB, err := database.Open()
		if err != nil {
			lg.Error(err)
			return c.JSON(http.StatusInternalServerError, response_types.Errors{})
		}

		cr := candidate_types.UpvoteResponse{}
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

				/* 候補日の存在チェック */
				_, err = candidate_logic.FetchCandidateByID(ctx, tx, cId, eId)
				if err != nil {
					return response_types.NewErrorSeed(http.StatusNotFound, "候補日が存在しません")
				}

				/* 投票の存在チェック */
				exists, err := candidate_logic.VoteExistsByUK(ctx, tx, accountId, cId)
				if err != nil {
					return err
				}
				if exists {
					err := candidate_logic.DeleteVote(ctx, tx, accountId, cId)
					if err != nil {
						return err
					}
				} else {
					/* 既に存在しないなら何もしない */
					lg.Debugf("no vote found. accountId: %d, candidateId: %d", accountId, cId)
				}

				ed, err := event_logic.GetDetail(ctx, tx, r.EventID)

				cr = candidate_types.UpvoteResponse{EventDetail: *ed}

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
		return c.JSON(http.StatusOK, cr)
	}
}
