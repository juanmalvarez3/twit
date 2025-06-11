package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorType string

const (
	ErrorTypeNotFound      ErrorType = "NOT_FOUND"
	ErrorTypeInvalidInput  ErrorType = "INVALID_INPUT"
	ErrorTypeUnauthorized  ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden     ErrorType = "FORBIDDEN"
	ErrorTypeAlreadyExists ErrorType = "ALREADY_EXISTS"
	ErrorTypeInternal      ErrorType = "INTERNAL"
	ErrorTypeUnavailable   ErrorType = "SERVICE_UNAVAILABLE"
	ErrorTypeInternalError ErrorType = "INTERNAL_ERROR"
)

type AppError struct {
	Type       ErrorType
	Message    string
	Cause      error
	StatusCode int
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (causa: %s)", e.Type, e.Message, e.Cause.Error())
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

func NewAppError(errType ErrorType, message string, cause error) *AppError {
	return &AppError{
		Type:       errType,
		Message:    message,
		Cause:      cause,
		StatusCode: mapErrorTypeToStatusCode(errType),
	}
}

func NewNotFoundError(message string, cause error) *AppError {
	return NewAppError(ErrorTypeNotFound, message, cause)
}

func NewInvalidInputError(message string, cause error) *AppError {
	return NewAppError(ErrorTypeInvalidInput, message, cause)
}

func NewInternalServerError(message string, cause error) *AppError {
	return NewAppError(ErrorTypeInternalError, message, cause)
}

func NewUnauthorizedError(message string, cause error) *AppError {
	return NewAppError(ErrorTypeUnauthorized, message, cause)
}

func NewForbiddenError(message string, cause error) *AppError {
	return NewAppError(ErrorTypeForbidden, message, cause)
}

func NewAlreadyExistsError(message string, cause error) *AppError {
	return NewAppError(ErrorTypeAlreadyExists, message, cause)
}

func NewInternalError(message string, cause error) *AppError {
	return NewAppError(ErrorTypeInternal, message, cause)
}

func NewUnavailableError(message string, cause error) *AppError {
	return NewAppError(ErrorTypeUnavailable, message, cause)
}

func Wrap(err error, message string) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		if message != "" {
			appErr.Message = fmt.Sprintf("%s: %s", message, appErr.Message)
		}
		return appErr
	}

	return NewInternalError(message, err)
}

func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.StatusCode
	}

	return http.StatusInternalServerError
}

func mapErrorTypeToStatusCode(errType ErrorType) int {
	switch errType {
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeInvalidInput:
		return http.StatusBadRequest
	case ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case ErrorTypeForbidden:
		return http.StatusForbidden
	case ErrorTypeAlreadyExists:
		return http.StatusConflict
	case ErrorTypeUnavailable:
		return http.StatusServiceUnavailable
	case ErrorTypeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
