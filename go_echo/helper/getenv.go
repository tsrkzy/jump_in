package helper

import (
	"log"
	"os"
)

// MustGetenv 環境変数を取得。なかったら(==空文字)エラーにする
func MustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Warning: %s environment variable not set.", k)
	}
	return v
}
