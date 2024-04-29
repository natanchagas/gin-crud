package customerrors

import (
	"net/http"
)

type ErrorCode string

var (
	UserRequestError ErrorCode = "BAD_REQUEST"
	ResourceNotFound ErrorCode = "RESOURCE_NOT_FOUND"
	ApplicationError ErrorCode = "APPLICATION_ERROR"
	UnexpectedError  ErrorCode = "UNEXPECTED_ERROR"
)

var (
	BadRequest = newError("something is wrong within your request", http.StatusBadRequest, UserRequestError)
	NotFound   = newError("resource not found", http.StatusNotFound, ResourceNotFound)
	Internal   = newError("application internal error", http.StatusInternalServerError, ApplicationError)
	Unexpected = newError("unexpected error", http.StatusInternalServerError, UnexpectedError)
)

type Error struct {
	StatusCode int
	ErrorCode  ErrorCode
	Message    string
}

func newError(message string, statusCode int, errorCode ErrorCode) Error {
	return Error{
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    message,
	}
}

func (e Error) Error() string {
	return e.Message
}
