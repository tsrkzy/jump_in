package response

type Errors struct {
	Errors []Error `json:"errors"`
}

func (s *Errors) Add(e Error) Errors {
	s.Errors = append(s.Errors, e)
	return *s
}

// ErrorGen
// 1つだけ Error を持つ Errors を返す時に便利なやつ
// DB接続エラーとかの toC で出せない無言エラーを生成するときは Errors{} で
func ErrorGen(msg string) Errors {
	s := Errors{}
	s.Add(Error{Msg: msg})
	return s
}

type Error struct {
	Msg string `json:"msg"`
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
