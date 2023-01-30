// Package lg
// "github.com/labstack/gommon/log" のラッパ
// ミドルウェアが呼び出されるまでは使用できない。 -> server.go
// もしそのタイミングでロギングをしたい場合は、組み込みの log パッケージを使用すること
package lg

import (
	"errors"
	"fmt"
	"github.com/tsrkzy/jump_in/helper/cx"
)

func Debug(i ...interface{}) {
	cxExist()
	//cx.Cx.Echo().Logger.Debug(i...)
	fmt.Print("Debug ")
	fmt.Println(i...)
}
func Debugf(format string, args ...interface{}) {
	cxExist()
	//cx.Cx.Echo().Logger.Debugf(format, args...)
	fmt.Print("Debug ")
	fmt.Printf(format, args...)
}
func Info(i ...interface{}) {
	cxExist()
	//cx.Cx.Echo().Logger.Info(i...)
	fmt.Print("Info ")
	fmt.Println(i...)
}
func Infof(format string, args ...interface{}) {
	cxExist()
	//cx.Cx.Echo().Logger.Infof(format, args...)
	fmt.Print("Info ")
	fmt.Printf(format, args...)
}

func Warn(i ...interface{}) {
	cxExist()
	//cx.Cx.Echo().Logger.Warn(i...)
	fmt.Print("Warn ")
	fmt.Println(i...)
}
func Warnf(format string, args ...interface{}) {
	cxExist()
	//cx.Cx.Echo().Logger.Warnf(format, args...)
	fmt.Print("Warn ")
	fmt.Printf(format, args...)
}

func Error(i ...interface{}) {
	cxExist()
	//cx.Cx.Echo().Logger.Error(i...)
	fmt.Print("Error ")
	fmt.Println(i...)
}
func Errorf(format string, args ...interface{}) {
	cxExist()
	//cx.Cx.Echo().Logger.Errorf(format, args...)
	fmt.Print("Error ")
	fmt.Printf(format, args...)
}

// 注意: FatalとPanicはその時点でPanicする == ログ吐いて死ぬ
//
// Fatal → その場で何もせずPanic
// Panic → JSONのHTTPレスポンスを返してからPanic

func Fatal(i ...interface{}) {
	cxExist()
	//cx.Cx.Echo().Logger.Fatal(i...)
	fmt.Print("Fatal ")
	fmt.Println(i...)
}
func Fatalf(format string, args ...interface{}) {
	cxExist()
	//cx.Cx.Echo().Logger.Fatalf(format, args...)
	fmt.Print("Fatal ")
	fmt.Printf(format, args...)
}
func Panic(i ...interface{}) {
	cxExist()
	//cx.Cx.Echo().Logger.Panic(i...)
	fmt.Print("Panic ")
	fmt.Println(i...)
}
func Panicf(format string, args ...interface{}) {
	cxExist()
	//cx.Cx.Echo().Logger.Panicf(format, args...)
	fmt.Print("Panic ")
	fmt.Printf(format, args...)
}

func cxExist() {
	if cx.Cx == nil {
		// ランタイムエラーだと分かりづらいのでここでpanicする
		panic(errors.New("do not call func in lg package before first middleware call"))
	}
}
