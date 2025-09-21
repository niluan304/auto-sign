package tieba

import (
	"fmt"
)

type Error struct {
	Code string `json:"error_code"`
	Msg  string `json:"error_msg"`
}

// 重复签到.
const codeSignRepeat = "160002"

// Error implements error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("code: %s, msg: %s", e.Code, e.Msg)
}
