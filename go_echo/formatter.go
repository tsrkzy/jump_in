package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"os"
	"strings"
)

/* アプリケーション内部で使用するロガーのフォーマッタ
 * @REFS https://ken-aio.github.io/post/2019/02/06/golang-echo-middleware/#logger%E3%83%9F%E3%83%89%E3%83%AB%E3%82%A6%E3%82%A7%E3%82%A2%E3%82%92%E3%82%AB%E3%82%B9%E3%82%BF%E3%83%9E%E3%82%A4%E3%82%BA%E3%81%99%E3%82%8B
 */
func internalLoggerFormatter(e *echo.Echo) {
	e.Logger.SetHeader("${id} [${level}] ")
	e.Logger.SetLevel(log.DEBUG)
}

/* requestLoggerFormatter
 * HTTPリクエスト用ロガーのカスタムフォーマッタ
 */
func requestLoggerFormatter() echo.MiddlewareFunc {
	logger := middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: formatter(),
		Output: os.Stdout,
	})
	return logger
}

/* formatLTSV */
func formatLTSV() string {
	lines := []string{
		"time:${time_rfc3339}",
		"host:${remote_ip}",
		"forwardedfor:${header:x-forwarded-for}",
		"req:-",
		"status:${status}",
		"method:${method}",
		"uri:${uri}",
		"size:${bytes_out}",
		"referer:${referer}",
		"ua:${user_agent}",
		"reqtime_ns:${latency}",
		"cache:-",
		"runtime:-",
		"apptime:-",
		"vhost:${host}",
		"reqtime_human:${latency_human}",
		"x-request-id:${id}",
		"host:${host}"}
	format := strings.Join(lines, "\t") + "\n"
	return format
}

func formatter() string {
	return `
- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - 
${method} ${host}${uri}
 - ${status}
 - size:${bytes_out}
 - reqtime_human:${latency_human}
 - referer:${referer}
 - ua:${user_agent}
================================================================================
`
}
