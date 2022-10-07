// Package lg
// "github.com/labstack/gommon/log" のラッパ
// ミドルウェアが呼び出されるまでは使用できない。 -> server.go
// もしそのタイミングでロギングをしたい場合は、組み込みの log パッケージを使用すること
package lg

import (
	"errors"
	"github.com/tsrkzy/jump_in/cx"
)

func Debug(i ...interface{}) {
	cxExist()
	cx.Cx.Echo().Logger.Debug(i...)
}
func Debugf(format string, args ...interface{}) {
	cxExist()
	cx.Cx.Echo().Logger.Debugf(format, args...)
}
func Info(i ...interface{}) {
	cxExist()
	cx.Cx.Echo().Logger.Info(i...)
}
func Infof(format string, args ...interface{}) {
	cxExist()
	cx.Cx.Echo().Logger.Infof(format, args...)
}

func Warn(i ...interface{}) {
	cxExist()
	cx.Cx.Echo().Logger.Warn(i...)
}
func Warnf(format string, args ...interface{}) {
	cxExist()
	cx.Cx.Echo().Logger.Warnf(format, args...)
}

func Error(i ...interface{}) {
	cxExist()
	cx.Cx.Echo().Logger.Error(i...)
}
func Errorf(format string, args ...interface{}) {
	cxExist()
	cx.Cx.Echo().Logger.Errorf(format, args...)
}

// 注意: FatalとPanicはその時点でPanicする == ログ吐いて死ぬ
//
// Fatal → その場で何もせずPanic
// Panic → JSONのHTTPレスポンスを返してからPanic

func Fatal(i ...interface{}) {
	cxExist()
	cx.Cx.Echo().Logger.Fatal(i...)
}
func Fatalf(format string, args ...interface{}) {
	cxExist()
	cx.Cx.Echo().Logger.Fatalf(format, args...)
}
func Panic(i ...interface{}) {
	cxExist()
	cx.Cx.Echo().Logger.Panic(i...)
}
func Panicf(format string, args ...interface{}) {
	cxExist()
	cx.Cx.Echo().Logger.Panicf(format, args...)
}

func cxExist() {
	if cx.Cx == nil {
		// ランタイムエラーだと分かりづらいのでここでpanicする
		panic(errors.New("do not call func in lg package before first middleware call"))
	}
}
