// Derror = Dave Error

package derror

import (
	"context"
	"fmt"
)

// Error represents the derror that implements the Go Error interface
type Error struct {
	Context   context.Context `json:"context"`
	Code      Code            `json:"code"`
	Type      Type            `json:"type"`
	Message   string          `json:"message"`
	Err       error           `json:"error"`
	Retryable bool            `json:"retryable"`
}

// New creates a new Derror that is not retryable
func New(ctx context.Context, code Code, errType Type, message string, err error) *Error {
	return &Error{
		Context:   ctx,
		Type:      errType,
		Code:      code,
		Message:   message,
		Retryable: false,
		Err:       err,
	}
}

// NewRetryable creates a new Derror with the ability to identify as retryable
func NewRetryable(ctx context.Context, code Code, errType Type, message string, err error) *Error {
	return &Error{
		Context:   ctx,
		Type:      errType,
		Code:      code,
		Message:   message,
		Retryable: true,
		Err:       err,
	}
}

// Error returns the string representation of the error
func (d *Error) Error() string {
	return fmt.Sprintf("code=%d, type=%s, message=%s, err=%v", d.Code, d.Type, d.Message, d.Err)
}

// IsRetryable returns whether the derror is retryable
func (d *Error) IsRetryable() bool {
	return d.Retryable
}

// GetContext returns the derror's context
func (d *Error) GetContext() context.Context {
	return d.Context
}
