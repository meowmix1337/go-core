package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HttpClient interface {
	Request(ctx context.Context, method string, endpoint string, payload interface{}, queryParams map[string]string) (*http.Response, error)
}

// httpClient implements the HttpClient interface
// takes in a base URL (https://dog.ceo) of the API and the API Prefix (v1/dog)
type httpClient struct {
	BaseURL   string
	APIPrefix string
}

// Request will make a request out the specified API endpoint
func (c *httpClient) Request(ctx context.Context, method string, endpoint string, payload interface{}, queryParams map[string]string) (*http.Response, error) {

	url := c.BaseURL + c.APIPrefix + endpoint

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if queryParams != nil {
		q := req.URL.Query()
		for k, v := range queryParams {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	if payload != nil {
		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
		req.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
		req.ContentLength = int64(len(jsonBytes))
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	// Check the response status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp, fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return resp, nil
}
