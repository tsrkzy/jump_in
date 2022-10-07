package validate

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/tsrkzy/jump_in/response"
)

// ErrorIntoJson validation.Validate の返す error を、
// response.Errors へ変換して返す
func ErrorIntoJson(e error) response.Errors {
	errs := e.(validation.Errors)

	vErr := response.Errors{}
	for _, err := range errs {
		m := fmt.Sprintf("%s", err)
		vErr.Add(response.Error{Msg: m})
	}
	return vErr
}
