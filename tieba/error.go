package tieba

import (
	"encoding/json"
	"fmt"
)

type Error struct {
	code string
	msg  string
}

func (e *Error) UnmarshalJSON(data []byte) error {
	var e1 struct {
		ErrorCode *string `json:"error_code"`
		ErrorMsg  *string `json:"error_msg"`

		No    int     `json:"no"`
		Error *string `json:"error"`
	}

	err := json.Unmarshal(data, &e1)
	if err != nil {
		return err
	}

	if e1.ErrorCode != nil {
		e.code = *e1.ErrorCode
	}
	if e1.ErrorMsg != nil {
		e.msg = *e1.ErrorMsg
	}

	if e1.No != 0 {
		e.code = Itoa(e1.No)
	}
	if e1.Error != nil {
		e.msg = *e1.Error
	}

	return nil
}

// CodeSignRepeat 重复签到.
const CodeSignRepeat = "160002"

// Error implements error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("code: %s, msg: %s", e.code, e.msg)
}

func (e *Error) Code() string {
	if e.code == "0" {
		return ""
	}

	return e.code
}

func (e *Error) Msg() string {
	return e.msg
}
