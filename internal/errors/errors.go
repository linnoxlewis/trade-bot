package errors

var (
	validation     = 1
	badRequest     = 400
	internal       = 500
	unauthorized   = 401
	accessDenied   = 403
	badGateway     = 502
	tooManyRequest = 429
	notFound       = 404
)

type Error struct {
	code int
	msg  string
	err  error
}

func New(code int, err string) Error {
	return Error{
		code: code,
		msg:  err,
	}
}

func (e Error) Error() string {
	return e.msg
}

func (e Error) GetCode() int {
	return e.code
}

func BadRequestError(err string) Error {
	return Error{
		code: badRequest,
		msg:  err,
		err:  nil,
	}
}

func WrapBadRequestError(err error) Error {
	return Error{
		err: err,
		msg: err.Error(),
	}
}

func (e Error) IsBadRequestError() bool {
	return e.GetCode() == badRequest
}

func InternalServerError(err error) Error {
	return Error{
		code: internal,
		msg:  "InternalServerError",
		err:  err,
	}
}

func (e Error) IsInternalServerError() bool {
	return e.GetCode() == internal
}

func ValidationError(err error) Error {
	return Error{
		code: validation,
		msg:  err.Error(),
		err:  err,
	}
}

func (e Error) IsValidationError() bool {
	return e.GetCode() == validation
}
