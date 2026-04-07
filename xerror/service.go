package xerror

import (
	"errors"
	"fmt"
)

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

func Extend(err IError, format string, args ...any) IError {
	return New(err.Code(), err.Message()+": "+fmt.Sprintf(format, args...))
}

func Error(err error, expects ...IError) bool {
	if err == nil {
		return false
	}
	var xerr IError
	ok := errors.As(err, &xerr)
	if !ok {
		return true
	}
	if xerr == nil {
		return false
	}
	for _, expect := range expects {
		if xerr.Code() == expect.Code() {
			return false
		}
	}
	return xerr.Code() != ErrSuccess.Code()
}
