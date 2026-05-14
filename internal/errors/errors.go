package errors

import "fmt"

// APIError represents a standardized API error response.
// @name APIError
type APIError struct {
	Code    string `json:"code" example:"INVALID_INPUT"`
	Message string `json:"message" example:"CPU usage cannot exceed CPU request for this analysis"`
}

// ErrorResponse wraps the APIError for the response body.
// @name ErrorResponse
type ErrorResponse struct {
	Error APIError `json:"error"`
}

// Common error codes
const (
	CodeInvalidInput   = "INVALID_INPUT"
	CodeInternalError  = "INTERNAL_ERROR"
	CodeNotFound       = "NOT_FOUND"
	CodeInvalidRequest = "INVALID_REQUEST"
)

// NewAPIError creates a new APIError.
func NewAPIError(code, message string) APIError {
	return APIError{
		Code:    code,
		Message: message,
	}
}

// Error implements the error interface.
func (e APIError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}
