package main

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/tsrkzy/jump_in/handler/account_handler"
	"github.com/tsrkzy/jump_in/handler/authenticate_handler"
	"github.com/tsrkzy/jump_in/handler/candidate_handler"
	"github.com/tsrkzy/jump_in/handler/event_handler"
	"github.com/tsrkzy/jump_in/helper"
	"github.com/tsrkzy/jump_in/helper/cx"
	"github.com/tsrkzy/jump_in/helper/validate"
	"net/http"
	"os"
	"time"
)

func main() {
	/* JST (GMT+9:00) */
	JST := time.FixedZone("Local", +9*60*60)
	time.Local = JST
	log.Info("started.")

	/* port listening */
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("defaulting to port %s", port)

	e := echo.New()

	/* バリデータ */
	e.Validator = &validate.CustomValidator{}

	/* middleware */
	e.Use(cx.CxWrapper())
	/*
	 * https://ken-aio.github.io/post/2019/02/06/golang-echo-middleware/
	 * https://qiita.com/rmanzoku/items/ddfde6d097443e634f22
	 * Panicした際に自動的に再起動する (systemdがあるので不要？)
	 */
	e.Use(middleware.Recover())

	/*
	 * HTTPリクエストにIDを振ってログに吐く
	 * @REFS https://echo.labstack.com/middleware/request-id/
	 */
	e.Use(middleware.RequestID())
	e.Use(requestLoggerFormatter())

	/* アプリケーション内部用ロガーのフォーマッタ */
	internalLoggerFormatter(e)

	/* セッションストア */
	sName := helper.MustGetenv("SESSION_NAME")
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(sName))))

	/* ************
	 * routing    *
	 * ************/

	/* --- 認証不要 --- */

	/* - ログインのためのAPI */
	e.POST("/api/authenticate", authenticate_handler.Authenticate())
	e.GET("/api/ml/:uri_hash", authenticate_handler.MagicLink())
	/* ログアウトはログオフ中でも可能 */
	e.GET("/api/logout", authenticate_handler.Logout())

	/* --- 認証必要 --- */

	e.GET("/api/status", authenticate_handler.Au(authenticate_handler.Status()))

	/* アカウント情報 */
	e.GET("/api/whoami", authenticate_handler.Au(authenticate_handler.WhoAmI()))
	e.POST("/api/account/name/update", authenticate_handler.Au(account_handler.UpdateName()))

	/* イベント */
	e.GET("/api/event/list", authenticate_handler.Au(event_handler.List()))
	e.GET("/api/event/detail", authenticate_handler.Au(event_handler.Detail()))
	e.POST("/api/event/create", authenticate_handler.Au(event_handler.Create()))
	e.POST("/api/event/name/update", authenticate_handler.Au(event_handler.UpdateName()))
	e.POST("/api/event/description/update", authenticate_handler.Au(event_handler.UpdateDescription()))
	e.POST("/api/event/open/update", authenticate_handler.Au(event_handler.UpdateOpen()))

	/* 候補日 */
	e.POST("/api/candidate/create", authenticate_handler.Au(candidate_handler.Create()))
	e.POST("/api/candidate/delete", authenticate_handler.Au(candidate_handler.Delete()))

	/* 投票 */
	e.POST("/api/vote/create", authenticate_handler.Au(candidate_handler.Upvote()))
	e.POST("/api/vote/delete", authenticate_handler.Au(candidate_handler.Downvote()))

	/* 参加/離脱 */
	e.POST("/api/event/attend", authenticate_handler.Au(event_handler.Attend()))
	e.POST("/api/event/leave", authenticate_handler.Au(event_handler.Leave()))

	/* start listening */
	if err := e.Start(":" + port); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
