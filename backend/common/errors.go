package common

import (
	"fmt"
	"net/http"
)

// AppError represents application-specific errors
type AppError struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
	StatusCode int    `json:"-"`
}

// Error implements the error interface
func (e AppError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

// NewAppError creates a new application error
func NewAppError(code int, message, details string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    details,
		StatusCode: http.StatusBadRequest,
	}
}

// NewValidationError creates a validation error
func NewValidationError(message string) *AppError {
	return &AppError{
		Code:       1001,
		Message:    "Validation Error",
		Details:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// NewNotFoundError creates a not found error
func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Code:       1002,
		Message:    "Not Found",
		Details:    fmt.Sprintf("%s not found", resource),
		StatusCode: http.StatusNotFound,
	}
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:       1003,
		Message:    "Unauthorized",
		Details:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

// NewInternalError creates an internal server error
func NewInternalError(message string) *AppError {
	return &AppError{
		Code:       1004,
		Message:    "Internal Server Error",
		Details:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

// NewServiceUnavailableError creates a service unavailable error
func NewServiceUnavailableError(service string) *AppError {
	return &AppError{
		Code:       1005,
		Message:    "Service Unavailable",
		Details:    fmt.Sprintf("%s service is currently unavailable", service),
		StatusCode: http.StatusServiceUnavailable,
	}
}

// Common error codes
const (
	ErrCodeValidation        = 1001
	ErrCodeNotFound         = 1002
	ErrCodeUnauthorized     = 1003
	ErrCodeInternal         = 1004
	ErrCodeServiceUnavailable = 1005
	ErrCodeOTPExpired       = 1006
	ErrCodeOTPInvalid       = 1007
	ErrCodeMaxAttempts      = 1008
	ErrCodeRateLimit        = 1009
) 