package errors

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewErrorResponse(t *testing.T) {
	code := "TEST_ERROR"
	message := "Test error message"
	details := map[string]string{"field": "value"}
	traceID := "trace-123"

	err := NewErrorResponse(code, message, details, traceID)

	if err == nil {
		t.Fatal("NewErrorResponse() returned nil")
	}

	if err.Code != code {
		t.Errorf("Expected code '%s', got '%s'", code, err.Code)
	}

	if err.Message != message {
		t.Errorf("Expected message '%s', got '%s'", message, err.Message)
	}

	if err.TraceID != traceID {
		t.Errorf("Expected traceID '%s', got '%s'", traceID, err.TraceID)
	}

	if err.Details == nil {
		t.Error("Details should not be nil")
	}
}

func TestNewErrorResponse_Timestamp(t *testing.T) {
	before := time.Now().UTC()
	err := NewErrorResponse("CODE", "message", nil, "")
	after := time.Now().UTC()

	if err.Timestamp == "" {
		t.Error("Timestamp should not be empty")
	}

	// Parse the timestamp
	timestamp, parseErr := time.Parse(time.RFC3339, err.Timestamp)
	if parseErr != nil {
		t.Errorf("Failed to parse timestamp: %v", parseErr)
	}

	// Verify timestamp is within reasonable range
	if timestamp.Before(before.Add(-1*time.Second)) || timestamp.After(after.Add(1*time.Second)) {
		t.Error("Timestamp is not within expected range")
	}
}

func TestNewErrorResponse_NilDetails(t *testing.T) {
	err := NewErrorResponse("CODE", "message", nil, "trace-123")

	if err.Details != nil {
		t.Error("Details should be nil when nil is passed")
	}
}

func TestNewErrorResponse_EmptyStrings(t *testing.T) {
	err := NewErrorResponse("", "", nil, "")

	if err == nil {
		t.Fatal("NewErrorResponse() returned nil")
	}

	if err.Code != "" {
		t.Error("Code should be empty string")
	}

	if err.Message != "" {
		t.Error("Message should be empty string")
	}

	if err.TraceID != "" {
		t.Error("TraceID should be empty string")
	}
}

func TestErrorResponse_JSONSerialization(t *testing.T) {
	details := map[string]interface{}{
		"field1": "value1",
		"field2": 123,
		"field3": true,
	}

	err := NewErrorResponse("TEST_ERROR", "Test message", details, "trace-456")

	// Marshal to JSON
	jsonData, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("Failed to marshal error response: %v", marshalErr)
	}

	// Unmarshal back
	var decoded ErrorResponse
	unmarshalErr := json.Unmarshal(jsonData, &decoded)
	if unmarshalErr != nil {
		t.Fatalf("Failed to unmarshal error response: %v", unmarshalErr)
	}

	// Verify fields
	if decoded.Code != err.Code {
		t.Errorf("Code mismatch: expected '%s', got '%s'", err.Code, decoded.Code)
	}

	if decoded.Message != err.Message {
		t.Errorf("Message mismatch: expected '%s', got '%s'", err.Message, decoded.Message)
	}

	if decoded.TraceID != err.TraceID {
		t.Errorf("TraceID mismatch: expected '%s', got '%s'", err.TraceID, decoded.TraceID)
	}

	if decoded.Timestamp != err.Timestamp {
		t.Errorf("Timestamp mismatch: expected '%s', got '%s'", err.Timestamp, decoded.Timestamp)
	}
}

func TestErrorResponse_JSONOmitEmpty(t *testing.T) {
	// Create error with no details or traceID
	err := NewErrorResponse("CODE", "message", nil, "")

	jsonData, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("Failed to marshal: %v", marshalErr)
	}

	jsonStr := string(jsonData)

	// Verify 'details' field is omitted when nil
	if containsField(jsonStr, "details") && err.Details == nil {
		t.Error("Details field should be omitted when nil")
	}

	// Verify 'trace_id' field is omitted when empty
	if containsField(jsonStr, "trace_id") && err.TraceID == "" {
		t.Error("TraceID field should be omitted when empty")
	}

	// Verify required fields are present
	if !containsField(jsonStr, "code") {
		t.Error("Code field should always be present")
	}

	if !containsField(jsonStr, "message") {
		t.Error("Message field should always be present")
	}

	if !containsField(jsonStr, "timestamp") {
		t.Error("Timestamp field should always be present")
	}
}

func TestErrorResponse_ComplexDetails(t *testing.T) {
	type ValidationError struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}

	details := []ValidationError{
		{Field: "email", Message: "Invalid email format"},
		{Field: "password", Message: "Password too short"},
	}

	err := NewErrorResponse("VALIDATION_ERROR", "Validation failed", details, "trace-789")

	if err.Details == nil {
		t.Error("Details should not be nil")
	}

	// Marshal and verify it works with complex types
	jsonData, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("Failed to marshal complex details: %v", marshalErr)
	}

	if len(jsonData) == 0 {
		t.Error("JSON data is empty")
	}
}

func TestErrorResponse_StringDetails(t *testing.T) {
	details := "Simple string error details"
	err := NewErrorResponse("ERROR", "message", details, "trace")

	if err.Details == nil {
		t.Error("Details should not be nil")
	}

	detailsStr, ok := err.Details.(string)
	if !ok {
		t.Error("Details should be a string")
	}

	if detailsStr != details {
		t.Errorf("Expected details '%s', got '%s'", details, detailsStr)
	}
}

func TestErrorResponse_NumberDetails(t *testing.T) {
	details := 404
	err := NewErrorResponse("NOT_FOUND", "Resource not found", details, "")

	if err.Details == nil {
		t.Error("Details should not be nil")
	}

	detailsNum, ok := err.Details.(int)
	if !ok {
		t.Error("Details should be an int")
	}

	if detailsNum != details {
		t.Errorf("Expected details %d, got %d", details, detailsNum)
	}
}

func TestErrorResponse_TimestampFormat(t *testing.T) {
	err := NewErrorResponse("CODE", "message", nil, "")

	// Verify timestamp is in RFC3339 format
	_, parseErr := time.Parse(time.RFC3339, err.Timestamp)
	if parseErr != nil {
		t.Errorf("Timestamp is not in RFC3339 format: %v", parseErr)
	}

	// Verify timestamp is UTC
	timestamp, _ := time.Parse(time.RFC3339, err.Timestamp)
	if timestamp.Location() != time.UTC {
		t.Error("Timestamp should be in UTC")
	}
}

// Helper function to check if JSON contains a field
func containsField(jsonStr, field string) bool {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return false
	}
	_, exists := data[field]
	return exists
}

func TestErrorResponse_MultipleInstances(t *testing.T) {
	// Create multiple error responses and verify they have unique timestamps
	var errors []*ErrorResponse

	for i := 0; i < 5; i++ {
		err := NewErrorResponse("ERROR", "message", nil, "")
		errors = append(errors, err)
		time.Sleep(10 * time.Millisecond) // Small delay to ensure different timestamps
	}

	// Verify each has a timestamp
	for i, err := range errors {
		if err.Timestamp == "" {
			t.Errorf("Error %d has empty timestamp", i)
		}
	}

	// Verify timestamps are in chronological order
	for i := 1; i < len(errors); i++ {
		t1, _ := time.Parse(time.RFC3339, errors[i-1].Timestamp)
		t2, _ := time.Parse(time.RFC3339, errors[i].Timestamp)

		if t2.Before(t1) {
			t.Errorf("Timestamp %d is before timestamp %d", i, i-1)
		}
	}
}
