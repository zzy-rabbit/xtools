package xerror

import "fmt"

type IError interface {
	Code() int
	Message() string
	Error() string
}

type err struct {
	ErrCode    int    `json:"code"`
	ErrMessage string `json:"message"`
}

func New(code int, message string) IError {
	return &err{
		ErrCode:    code,
		ErrMessage: message,
	}
}

func (e *err) Code() int {
	return e.ErrCode
}

func (e *err) Message() string {
	return e.ErrMessage
}

func (e *err) Error() string {
	return fmt.Sprintf("%d: %s", e.ErrCode, e.ErrMessage)
}
