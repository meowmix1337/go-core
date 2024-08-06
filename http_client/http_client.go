package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/meowmix1337/go-core/derror"
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
		return nil, derror.New(ctx, derror.InternalServerCode, derror.InternalType, "failed to create new request", err)
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
			return nil, derror.New(ctx, derror.InternalServerCode, derror.InternalType, "failed to marshal the payload", err)
		}
		req.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
		req.ContentLength = int64(len(jsonBytes))
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, derror.New(ctx, derror.InternalServerCode, derror.InternalType, "failed to do request", err)
	}

	// Check the response status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp, derror.New(ctx, derror.InternalServerCode, derror.InternalType, "request response received a bad status code", errors.New("bad response code"))
	}

	return resp, nil
}
