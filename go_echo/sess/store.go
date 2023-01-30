package sess

import (
	"github.com/antonlindstrom/pgstore"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/tsrkzy/jump_in/database"
	"github.com/tsrkzy/jump_in/helper"
	"net/http"
)

var svNameChocochip string

// SvNameChocochip
// MLを要求したデバイスと、ログインを試みたデバイスが同一であることを検証するためのセッション変数名
// ログイン後は寿命を延長し、認証中であることを示すセッション変数へ昇格する
func SvNameChocochip() string {
	if svNameChocochip == "" {
		log.Debug("SvNameChocochip initial load")
		svNameChocochip = helper.MustGetenv("SESSION_VAR_NAME")
	}
	return svNameChocochip
}

func Open(c echo.Context, db *database.MyDB, callbackFn func(*sessions.Session) error) error {
	/* 引数として受け取ったDB接続とキーから、DBのhttp_sessionテーブル読み書き用構造体を作成 */
	storeKey := helper.MustGetenv("SESSION_STORE_SECURE_KEY")
	pgs, err := pgstore.NewPGStoreFromPool(db.GetDB(), []byte(storeKey))
	if err != nil {
		return err
	}
	defer pgs.Close()

	/* cookieからセッション取得 */
	sName := helper.MustGetenv("SESSION_NAME")
	ss, err := pgs.Get(c.Request(), sName)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	defer func() {
		if p := recover(); p != nil {
			log.Error(p)
			panic(p)
		}
		/* session.Saveのerrをそのまま外に流す */
		err = ss.Save(c.Request(), c.Response())
	}()

	err = callbackFn(ss)
	return err
}
