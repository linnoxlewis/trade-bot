package errors

import (
	"fmt"
	"net/http"
)

const (
	ValidationErrorCode = iota + 1
)

type Error struct {
	code int
	msg  string
	err  error
}

func New(code int, msg string) Error {
	return Error{
		code: code,
		msg:  msg,
	}
}

func (e Error) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.err)
	}
	return e.msg
}

func (e Error) GetCode() int {
	return e.code
}

func (e Error) Unwrap() error {
	return e.err
}

func BadRequestError(msg string) Error {
	return Error{
		code: http.StatusBadRequest,
		msg:  msg,
	}
}

func WrapBadRequestError(err error) Error {
	return Error{
		code: http.StatusBadRequest,
		msg:  err.Error(),
		err:  err,
	}
}

func (e Error) IsBadRequestError() bool {
	return e.GetCode() == http.StatusBadRequest
}

func InternalServerError(err error) Error {
	return Error{
		code: http.StatusInternalServerError,
		err:  err,
	}
}

func (e Error) IsInternalServerError() bool {
	return e.GetCode() == http.StatusInternalServerError
}

func ValidationError(msg string, err error) Error {
	return Error{
		code: ValidationErrorCode,
		msg:  msg,
		err:  err,
	}
}

func (e Error) IsValidationError() bool {
	return e.GetCode() == ValidationErrorCode
}
