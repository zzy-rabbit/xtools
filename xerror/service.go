package xerror

import "fmt"

type IError interface {
	Code() int
	Message() string
	Error() string
}

type err struct {
	code    int
	message string
}

func New(code int, message string) IError {
	return &err{
		code:    code,
		message: message,
	}
}

func (e *err) Code() int {
	return e.code
}

func (e *err) Message() string {
	return e.message
}

func (e *err) Error() string {
	return fmt.Sprintf("%d: %s", e.code, e.message)
}
