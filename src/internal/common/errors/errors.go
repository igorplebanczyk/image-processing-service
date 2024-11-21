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
	typ     ErrorType
	message string
}

func (e Error) Error() string {
	return e.message
}

func (e Error) Type() ErrorType {
	return e.typ
}

func New(message string) Error {
	return Error{
		typ:     Unknown,
		message: message,
	}
}

func NewInvalidInput(message string) Error {
	return Error{
		typ:     InvalidInput,
		message: message,
	}
}

func NewUnauthorized(message string) Error {
	return Error{
		typ:     Unauthorized,
		message: message,
	}
}

func NewForbidden(message string) Error {
	return Error{
		typ:     Forbidden,
		message: message,
	}
}

func NewInternal(message string) Error {
	return Error{
		typ:     Internal,
		message: message,
	}
}
