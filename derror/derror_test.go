package derror_test

import (
	"context"
	"errors"
	"testing"

	"github.com/meowmix1337/go-core/derror"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	err := derror.New(context.Background(), derror.InternalServerCode, derror.InternalType, "failed to do something", errors.New("some error"))

	assert.NotNil(t, err)
	assert.EqualError(t, err, "code=500, type=INTERNAL_ERROR, message=failed to do something, err=some error")
	assert.False(t, err.IsRetryable())
}

func TestNewRetryable(t *testing.T) {
	err := derror.NewRetryable(context.Background(), derror.InternalServerCode, derror.InternalType, "failed to do something", errors.New("some error"))

	assert.NotNil(t, err)
	assert.EqualError(t, err, "code=500, type=INTERNAL_ERROR, message=failed to do something, err=some error")
	assert.True(t, err.IsRetryable())
}

func TestWrap(t *testing.T) {
	parentErr := derror.New(context.Background(), derror.InternalServerCode, derror.InternalType, "this is the parent", errors.New("parent error"))
	wrappedErr := derror.New(context.Background(), derror.BadRequestCode, derror.BadRequestType, "this is the child", parentErr)

	assert.NotNil(t, wrappedErr)
	assert.EqualError(t, wrappedErr, "code=400, type=BAD_REQUEST, message=this is the child, err=code=500, type=INTERNAL_ERROR, message=this is the parent, err=parent error")
	assert.False(t, wrappedErr.IsRetryable())

	assert.EqualError(t, wrappedErr.Unwrap(), "code=500, type=INTERNAL_ERROR, message=this is the parent, err=parent error")
}
