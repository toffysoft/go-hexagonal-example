package errors

import "net/http"

type ErrorType string

func (e ErrorType) String() string {
	return string(e)
}

const (
	NotFound       ErrorType = "NOT_FOUND"
	InvalidInput   ErrorType = "INVALID_INPUT"
	InternalServer ErrorType = "INTERNAL_SERVER_ERROR"
	Unauthorized   ErrorType = "UNAUTHORIZED"
	Forbidden      ErrorType = "FORBIDDEN"
)

type AppError struct {
	Type    ErrorType
	Message string
}

func (e AppError) Error() string {
	return e.Message
}

func (e AppError) StatusCode() int {
	switch e.Type {
	case NotFound:
		return http.StatusNotFound
	case InvalidInput:
		return http.StatusBadRequest
	case Unauthorized:
		return http.StatusUnauthorized
	case Forbidden:
		return http.StatusForbidden
	case InternalServer:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func NewAppError(errorType ErrorType, message string) AppError {
	return AppError{
		Type:    errorType,
		Message: message,
	}
}

func NewNotFoundError(message string) AppError {
	return NewAppError(NotFound, message)
}

func NewInvalidInputError(message string) AppError {
	return NewAppError(InvalidInput, message)
}

func NewInternalServerError(message string) AppError {
	return NewAppError(InternalServer, message)
}

func NewUnauthorizedError(message string) AppError {
	return NewAppError(Unauthorized, message)
}

func NewForbiddenError(message string) AppError {
	return NewAppError(Forbidden, message)
}
