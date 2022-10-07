package validate

import validation "github.com/go-ozzo/ozzo-validation"

// CustomValidator
// Validate メソッドで任意の構造体を引数として受け取る
// validation.Validatable に埋め込まれた、バリデーション可能なもの
// ( ライブラリ内で Validate が実装されているデータ型 )
// の場合のみバリデーションを実施し、違反がある場合のみ error を返す
type CustomValidator struct {
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if c, ok := i.(validation.Validatable); ok {
		return c.Validate()
	}
	return nil
}
