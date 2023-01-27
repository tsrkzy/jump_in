package cx

import "github.com/labstack/echo/v4"

// Cx
//
// Cx はスレッド(==goroutine)セーフではない
// @REFS server.Serve /src/net/http/server.go
//
// ロガー用にcontext.Contextをラップして使いまわしている
// server.Serveはリクエストが来るたびに goroutine を生成するため、
// Cx はそのリクエストに対するcontextである保証がない
// ロガー以外の用途( クエリパラメータの取得など )に使用しないこと
var Cx *CustomContext

type CustomContext struct {
	echo.Context
}

func CxWrapper() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			// 最初に呼び出されるMiddlewareで、contextをラップして次のMiddleWareに渡す
			// @REFS https://echo.labstack.com/guide/context/#extending-context
			Cx = &CustomContext{c}
			return next(Cx)
		}
	}
}
