package http_util

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJSONResponse(t *testing.T) {
	// Create a response recorder to capture the HTTP response
	recorder := httptest.NewRecorder()

	// Example data to be returned as JSON
	data := map[string]string{"message": "Hello, World!"}

	// Call the JSONResponse function with the recorder as the response writer
	JSONResponse(recorder, http.StatusOK, data)

	// Check the HTTP status code
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	// Check the Content-Type header
	contentType := recorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// Check the response body
	expectedBody := `{"message":"Hello, World!"}`

	actualBody := strings.TrimSpace(recorder.Body.String())
	if actualBody != expectedBody {
		t.Errorf("Expected response body %s, got %s", expectedBody, actualBody)
	}
}
