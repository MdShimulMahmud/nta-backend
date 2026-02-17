package utils

import (
	"encoding/json"
	"fmt"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// SuccessResponse creates a success response
func SuccessResponse(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

// ErrorResponse creates an error response
func ErrorResponse(err error) Response {
	return Response{
		Success: false,
		Error:   err.Error(),
	}
}

// MessageResponse creates a message response
func MessageResponse(message string) Response {
	return Response{
		Success: true,
		Message: message,
	}
}

// PrettyPrint prints JSON in a pretty format
func PrettyPrint(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("%+v", v)
	}
	return string(b)
}
