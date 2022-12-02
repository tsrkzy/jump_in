package helper

import (
	"github.com/friendsofgo/errors"
	"strconv"
)

func StrToID(numStr string) (int64, error) {
	i, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, err
	}
	if i < 0 {
		return 0, errors.New("正の整数で指定してください")
	}
	return int64(i), nil
}
