package apiclient_test

import (
	"net/http"
	"testing"

	"github.com/meowmix1337/go-core/apiclient"
	"github.com/stretchr/testify/assert"
)

func TestRequestMethod_String(t *testing.T) {
	tests := []struct {
		name     string
		method   apiclient.RequestMethod
		expected string
	}{
		{"GET method", apiclient.MethodGet, http.MethodGet},
		{"POST method", apiclient.MethodPost, http.MethodPost},
		{"PUT method", apiclient.MethodPut, http.MethodPut},
		{"DELETE method", apiclient.MethodDelete, http.MethodDelete},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.method.String())
		})
	}
}
