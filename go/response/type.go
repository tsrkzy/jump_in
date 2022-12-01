package response

import "fmt"

type Errors struct {
	Errors []ErrorSeed `json:"errors"`
}

func (s *Errors) Add(e ErrorSeed) Errors {
	s.Errors = append(s.Errors, e)
	return *s
}

// ErrorGen
// 1つだけ ErrorSeed を持つ Errors を返す時に便利なやつ
// DB接続エラーとかの toC で出せない無言エラーを生成するときは Errors{} で
func ErrorGen(msg string) Errors {
	s := Errors{}
	s.Add(ErrorSeed{Msg: msg})
	return s
}

// ErrorSeed ロジックの中からステータスコードを持った error を返す用
// NewErrorSeed から作成する想定
type ErrorSeed struct {
	Code int    `json:"-"`
	Msg  string `json:"msg"`
}

// Error error の interface の実装
func (e ErrorSeed) Error() string {
	return fmt.Sprintf("%s", e.Msg)
}

// ErrorSeed を作成して返す
func NewErrorSeed(code int, msg string) ErrorSeed {
	return ErrorSeed{code, msg}
}

// OK
// 処理自体は正常に終わったけど特に返すものがないよって時の JSON レスポンス用
// OK.Ok には true しか入らない想定 Ok() で作ると楽
type OK struct {
	Ok bool `json:"ok"`
}

// Ok
// @SEE OK
func Ok() OK {
	return OK{Ok: true}
}
