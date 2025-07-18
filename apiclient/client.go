package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// HTTPClient defines the interface for making HTTP requests
type HTTPClient interface {
	Do(ctx context.Context, req Request, result interface{}) error
	Get(ctx context.Context, endpoint string, queryParams map[string]string, headers map[string]string, result interface{}) error
	Post(ctx context.Context, endpoint string, body interface{}, headers map[string]string, result interface{}) error
	Put(ctx context.Context, endpoint string, body interface{}, headers map[string]string, result interface{}) error
	Delete(ctx context.Context, endpoint string, headers map[string]string, result interface{}) error
}

// APIClient is a simple HTTP client wrapper for making API requests
type APIClient struct {
	baseURL        string
	httpClient     *http.Client
	defaultHeaders map[string]string
}

// Option defines a function type for APIClient configuration options
type Option func(*APIClient)

func WithBaseURL(baseURL string) Option {
	return func(c *APIClient) {
		c.baseURL = baseURL
	}
}

// WithTimeout sets the timeout for the HTTP client
func WithTimeout(timeout time.Duration) Option {
	return func(c *APIClient) {
		c.httpClient.Timeout = timeout
	}
}

// WithDefaultHeader adds a default header to all requests
func WithAuthorization(token string) Option {
	return func(c *APIClient) {
		c.defaultHeaders["Authorization"] = "Bearer " + token
	}
}

// NewAPIClient creates a new API client with default configuration
func NewAPIClient(opts ...Option) *APIClient {
	client := &APIClient{
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
		defaultHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	// Apply all options
	for _, opt := range opts {
		opt(client)
	}

	return client
}

// Do performs an HTTP request and returns the response
func (c *APIClient) Do(ctx context.Context, req Request, result interface{}) error {
	// Construct full URL
	fullUrl := c.baseURL + req.Endpoint

	// Create request body if any
	var bodyReader io.Reader
	if req.Body != nil {
		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, req.Method.String(), fullUrl, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set default headers
	for k, v := range c.defaultHeaders {
		httpReq.Header.Set(k, v)
	}

	// Set request-specific headers
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	// Add query parameters if any
	if len(req.QueryParams) > 0 {
		q := url.Values{}
		for key, value := range req.QueryParams {
			q.Add(key, value)
		}
		httpReq.URL.RawQuery = q.Encode()
	}

	// Execute request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response if result pointer is provided
	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// Get is a convenience method for making GET requests
func (c *APIClient) Get(ctx context.Context, endpoint string, queryParams map[string]string, headers map[string]string, result interface{}) error {
	return c.Do(ctx, Request{
		Method:      MethodGet,
		Endpoint:    endpoint,
		QueryParams: queryParams,
		Headers:     headers,
	}, result)
}

// Post is a convenience method for making POST requests
func (c *APIClient) Post(ctx context.Context, endpoint string, body interface{}, headers map[string]string, result interface{}) error {
	return c.Do(ctx, Request{
		Method:   MethodPost,
		Endpoint: endpoint,
		Body:     body,
		Headers:  headers,
	}, result)
}

// Put is a convenience method for making PUT requests
func (c *APIClient) Put(ctx context.Context, endpoint string, body interface{}, headers map[string]string, result interface{}) error {
	return c.Do(ctx, Request{
		Method:   MethodPut,
		Endpoint: endpoint,
		Body:     body,
		Headers:  headers,
	}, result)
}

// Delete is a convenience method for making DELETE requests
func (c *APIClient) Delete(ctx context.Context, endpoint string, headers map[string]string, result interface{}) error {
	return c.Do(ctx, Request{
		Method:   MethodDelete,
		Endpoint: endpoint,
		Headers:  headers,
	}, result)
}
