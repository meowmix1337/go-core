package apiclient_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	apiclient "github.com/meowmix1337/go-core/api_client"
	"github.com/stretchr/testify/assert"
)

func TestDo_Success(t *testing.T) {
	// Mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/test-endpoint", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "value", r.URL.Query().Get("key"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "success"}`))
	}))
	defer ts.Close()

	client := apiclient.NewAPIClient(apiclient.WithBaseURL(ts.URL), apiclient.WithAuthorization("test-token"))

	var result map[string]string
	err := client.Do(context.Background(), apiclient.Request{
		Method:      apiclient.MethodGet,
		Endpoint:    "/test-endpoint",
		QueryParams: map[string]string{"key": "value"},
	}, &result)

	assert.NoError(t, err)
	assert.Equal(t, "success", result["message"])
}

func TestDo_Non2xxStatusCode(t *testing.T) {
	// Mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error": "bad request"}`))
	}))
	defer ts.Close()

	client := apiclient.NewAPIClient(apiclient.WithBaseURL(ts.URL))

	err := client.Do(context.Background(), apiclient.Request{
		Method:   apiclient.MethodGet,
		Endpoint: "/test-endpoint",
	}, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API request failed with status 400")
	assert.Contains(t, err.Error(), "bad request")
}

func TestDo_InvalidJSONResponse(t *testing.T) {
	// Mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`invalid-json`))
	}))
	defer ts.Close()

	client := apiclient.NewAPIClient(apiclient.WithBaseURL(ts.URL))

	var result map[string]string
	err := client.Do(context.Background(), apiclient.Request{
		Method:   apiclient.MethodGet,
		Endpoint: "/test-endpoint",
	}, &result)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode response")
}

func TestDo_NetworkError(t *testing.T) {
	client := apiclient.NewAPIClient(apiclient.WithBaseURL("http://invalid-url"))

	err := client.Do(context.Background(), apiclient.Request{
		Method:   apiclient.MethodGet,
		Endpoint: "/test-endpoint",
	}, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request failed")
}

func TestDo_RequestBody(t *testing.T) {
	// Mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"key":"value"}`, string(body))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "success"}`))
	}))
	defer ts.Close()

	client := apiclient.NewAPIClient(apiclient.WithBaseURL(ts.URL))

	var result map[string]string
	err := client.Do(context.Background(), apiclient.Request{
		Method:   apiclient.MethodPost,
		Endpoint: "/test-endpoint",
		Body:     map[string]string{"key": "value"},
	}, &result)

	assert.NoError(t, err)
	assert.Equal(t, "success", result["message"])
}

func TestDo_DefaultHeaders(t *testing.T) {
	// Mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := apiclient.NewAPIClient(apiclient.WithBaseURL(ts.URL))

	err := client.Do(context.Background(), apiclient.Request{
		Method:   apiclient.MethodGet,
		Endpoint: "/test-endpoint",
	}, nil)

	assert.NoError(t, err)
}

func TestGet(t *testing.T) {
	// Mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/test-endpoint", r.URL.Path)
		assert.Equal(t, "value", r.URL.Query().Get("key"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "success"}`))
	}))
	defer ts.Close()

	client := apiclient.NewAPIClient(apiclient.WithBaseURL(ts.URL))

	var result map[string]string
	err := client.Get(context.Background(), "/test-endpoint", map[string]string{"key": "value"}, nil, &result)

	assert.NoError(t, err)
	assert.Equal(t, "success", result["message"])
}

func TestPost(t *testing.T) {
	// Mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"key":"value"}`, string(body))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "created"}`))
	}))
	defer ts.Close()

	client := apiclient.NewAPIClient(apiclient.WithBaseURL(ts.URL))

	var result map[string]string
	err := client.Post(context.Background(), "/test-endpoint", map[string]string{"key": "value"}, nil, &result)

	assert.NoError(t, err)
	assert.Equal(t, "created", result["message"])
}

func TestPut(t *testing.T) {
	// Mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"key":"value"}`, string(body))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "updated"}`))
	}))
	defer ts.Close()

	client := apiclient.NewAPIClient(apiclient.WithBaseURL(ts.URL))

	var result map[string]string
	err := client.Put(context.Background(), "/test-endpoint", map[string]string{"key": "value"}, nil, &result)

	assert.NoError(t, err)
	assert.Equal(t, "updated", result["message"])
}

func TestDelete(t *testing.T) {
	// Mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/test-endpoint", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "deleted"}`))
	}))
	defer ts.Close()

	client := apiclient.NewAPIClient(apiclient.WithBaseURL(ts.URL))

	var result map[string]string
	err := client.Delete(context.Background(), "/test-endpoint", nil, &result)

	assert.NoError(t, err)
	assert.Equal(t, "deleted", result["message"])
}
