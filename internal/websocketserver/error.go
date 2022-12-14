package websocketserver

import "errors"

const (
	CodeAck              = iota
	CodeCommonError      = iota + 2000
	CodeRequestDataError = iota + 2001
	CodeInternalError
)

type Error struct {
	Code int
	error
}

func (err *Error) Error() string {
	return err.error.Error()
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:  code,
		error: errors.New(message),
	}
}
