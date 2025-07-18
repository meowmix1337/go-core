package apiclient

import "net/http"

type RequestMethod string

func (m RequestMethod) String() string {
	return string(m)
}

const (
	MethodGet    RequestMethod = http.MethodGet
	MethodPost   RequestMethod = http.MethodPost
	MethodPut    RequestMethod = http.MethodPut
	MethodDelete RequestMethod = http.MethodDelete
)
