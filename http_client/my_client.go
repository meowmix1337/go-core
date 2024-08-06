package httpclient

import (
	"context"
	"net/http"
)

// MyClient is a wrapper struct that implements specific HTTP methods
type MyClient struct {
	client HttpClient
}

func New(baseUrl, apiPrefix string) *MyClient {
	return &MyClient{
		client: &httpClient{
			BaseURL:   baseUrl,
			APIPrefix: apiPrefix,
		},
	}
}

func (c *MyClient) Get(ctx context.Context, endpoint string, queryParams map[string]string) (*http.Response, error) {
	return c.client.Request(ctx, http.MethodGet, endpoint, nil, queryParams)
}

func (c *MyClient) Post(ctx context.Context, endpoint string, payload interface{}) (*http.Response, error) {
	return c.client.Request(ctx, http.MethodPost, endpoint, payload, nil)
}

func (c *MyClient) Put(ctx context.Context, endpoint string, payload interface{}) (*http.Response, error) {
	return c.client.Request(ctx, http.MethodPut, endpoint, payload, nil)
}

func (c *MyClient) Delete(ctx context.Context, endpoint string, payload interface{}) (*http.Response, error) {
	return c.client.Request(ctx, http.MethodDelete, endpoint, payload, nil)
}
