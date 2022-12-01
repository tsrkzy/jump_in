package main

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/tsrkzy/jump_in/authenticate"
	"github.com/tsrkzy/jump_in/cx"
	"github.com/tsrkzy/jump_in/event"
	"github.com/tsrkzy/jump_in/helper"
	"github.com/tsrkzy/jump_in/response"
	"github.com/tsrkzy/jump_in/validate"
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

	/* routing */
	/* 認証不要 - ログインのためのAPI */
	e.POST("/api/authenticate", authenticate.Authenticate())
	e.GET("/api/ml/:uri_hash", authenticate.MagicLink())
	/* ログアウトはログオフ中でも可能 */
	e.GET("/api/logout", authenticate.Logout())

	/* ログインが必要 */
	e.GET("/api/status", au(authenticate.Status()))
	e.GET("/api/event/list", au(event.List()))
	e.GET("/api/event/detail", au(event.Detail()))
	e.POST("/api/event/create", au(event.Create()))
	e.POST("/api/event/attend", au(event.Attend()))

	/* start listening */
	if err := e.Start(":" + port); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// au
// API(HandlerFunc)の呼び出しにログイン認証をかける
// ログイン済みならそのまま Handler を処理
// 未ログインなら 401 Unauthorized で応答
func au(f echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := authenticate.Authorized(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, response.ErrorGen("ログインが必要です"))
		}

		return f(c)
	}
}
