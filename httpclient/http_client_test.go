package httpclient

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	// Create a test server to handle requests
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/data" {
			t.Errorf("expected URL path /api/data, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"data": "test data"}`))
		if err != nil {
			t.Errorf("error writing response: %v", err)
		}
	}))
	defer ts.Close()

	httpClient := httpClient{
		BaseURL:   ts.URL,
		APIPrefix: "/api",
	}

	// Define test cases
	tests := []struct {
		name        string
		method      string
		endpoint    string
		payload     interface{}
		queryParams map[string]string
		expected    interface{}
	}{
		{
			name:        "GET request with no payload or query params",
			method:      "GET",
			endpoint:    "/data",
			payload:     nil,
			queryParams: nil,
			expected:    map[string]interface{}{"data": "test data"},
		},
		{
			name:        "GET request with query params",
			method:      "GET",
			endpoint:    "/data",
			payload:     nil,
			queryParams: map[string]string{"string": "bar", "num": "1", "bool": "true"},
			expected:    map[string]interface{}{"data": "test data"},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := httpClient.Request(context.Background(), tt.method, tt.endpoint, tt.payload, tt.queryParams)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("expected no error reading response body, got %v", err)
			}

			var actual interface{}
			err = json.Unmarshal(body, &actual)
			if err != nil {
				t.Errorf("expected no error unmarshaling response, got %v", err)
			}
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected response body %v, got %v", tt.expected, actual)
			}
		})
	}
}

func TestPost(t *testing.T) {
	// Create a test server to handle requests
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/data" {
			t.Errorf("expected URL path /api/data, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte(`{"data": "test data"}`))
		if err != nil {
			t.Errorf("error writing response: %v", err)
		}
	}))
	defer ts.Close()

	httpClient := httpClient{
		BaseURL:   ts.URL,
		APIPrefix: "/api",
	}

	// Define test cases
	tests := []struct {
		name        string
		method      string
		endpoint    string
		payload     interface{}
		queryParams map[string]string
		expected    interface{}
	}{
		{
			name:        "POST request with payload",
			method:      "POST",
			endpoint:    "/data",
			payload:     map[string]interface{}{"foo": "bar", "num": 1},
			queryParams: nil,
			expected:    map[string]interface{}{"data": "test data"},
		},
		{
			name:        "POST request without payload",
			method:      "POST",
			endpoint:    "/data",
			payload:     nil,
			queryParams: nil,
			expected:    map[string]interface{}{"data": "test data"},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := httpClient.Request(context.Background(), tt.method, tt.endpoint, tt.payload, tt.queryParams)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if resp.StatusCode != http.StatusCreated {
				t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("expected no error reading response body, got %v", err)
			}

			var actual interface{}
			err = json.Unmarshal(body, &actual)
			if err != nil {
				t.Errorf("expected no error unmarshaling response, got %v", err)
			}
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected response body %v, got %v", tt.expected, actual)
			}
		})
	}
}

func TestRequest_ErrorStatusCode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error": "server error"}`))
		if err != nil {
			t.Errorf("error writing response: %v", err)
		}
	}))
	defer ts.Close()

	c := &httpClient{
		BaseURL: ts.URL,
	}
	resp, err := c.Request(context.Background(), "GET", "", nil, nil)
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	assert.NotEmpty(t, err)
	assert.Equal(t, "code=500, type=INTERNAL_ERROR, message=request response received a bad status code, err=bad response code", err.Error())
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
