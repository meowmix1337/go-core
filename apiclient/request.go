package apiclient

// Request contains parameters for making an API request
type Request struct {
	Method      RequestMethod
	Endpoint    string
	Body        interface{}
	QueryParams map[string]string
	Headers     map[string]string
}