package validate

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/tsrkzy/jump_in/types/response_types"
)

// ErrorIntoJson validation.Validate の返す error を、
// response.Errors へ変換して返す
func ErrorIntoJson(e error) response_types.Errors {
	errs := e.(validation.Errors)

	vErr := response_types.Errors{}
	for _, err := range errs {
		m := fmt.Sprintf("%s", err)
		vErr.Add(response_types.ErrorSeed{Msg: m})
	}
	return vErr
}
