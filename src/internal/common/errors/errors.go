package errors

type ErrorType string

const (
	InvalidInput ErrorType = "invalid_input"
	Unauthorized ErrorType = "unauthorized"
	Forbidden    ErrorType = "forbidden"
	Internal     ErrorType = "internal"
	Unknown      ErrorType = "unknown"
)

type Error struct {
	typ ErrorType
	msg string
}

func (e Error) Error() string {
	return e.msg
}

func (e Error) Type() ErrorType {
	return e.typ
}

func New(message string) Error {
	return Error{
		typ: Unknown,
		msg: message,
	}
}

func NewInvalidInput(message string) Error {
	return Error{
		typ: InvalidInput,
		msg: message,
	}
}

func NewUnauthorized(message string) Error {
	return Error{
		typ: Unauthorized,
		msg: message,
	}
}

func NewForbidden(message string) Error {
	return Error{
		typ: Forbidden,
		msg: message,
	}
}

func NewInternal(message string) Error {
	return Error{
		typ: Internal,
		msg: message,
	}
}
