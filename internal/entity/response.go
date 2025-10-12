package entity

import (
	"net/http"
)

// APIResponse represents the standard API response format
type APIResponse struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse creates a success response
func SuccessResponse(code int, message string, data interface{}) *APIResponse {
	return &APIResponse{
		Status:  "success",
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse creates an error response
func ErrorResponse(code int, message string) *APIResponse {
	return &APIResponse{
		Status:  "fail",
		Code:    code,
		Message: message,
		Data:    nil,
	}
}

// Common success responses
func OKResponse(message string, data interface{}) *APIResponse {
	return SuccessResponse(http.StatusOK, message, data)
}

func CreatedResponse(message string, data interface{}) *APIResponse {
	return SuccessResponse(http.StatusCreated, message, data)
}

// Common error responses
func BadRequestResponse(message string) *APIResponse {
	return ErrorResponse(http.StatusBadRequest, message)
}

func UnauthorizedResponse(message string) *APIResponse {
	return ErrorResponse(http.StatusUnauthorized, message)
}

func NotFoundResponse(message string) *APIResponse {
	return ErrorResponse(http.StatusNotFound, message)
}

func ConflictResponse(message string) *APIResponse {
	return ErrorResponse(http.StatusConflict, message)
}

func InternalServerErrorResponse(message string) *APIResponse {
	return ErrorResponse(http.StatusInternalServerError, message)
}
