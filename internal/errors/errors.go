package errors

import "time"

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details,omitempty"`
	TraceID   string      `json:"trace_id,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(code, message string, details interface{}, traceID string) *ErrorResponse {
	return &ErrorResponse{
		Code:      code,
		Message:   message,
		Details:   details,
		TraceID:   traceID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}
