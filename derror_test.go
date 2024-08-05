package derror_test

import (
	"context"
	"errors"
	"testing"

	"github.com/meowmix1337/derror"
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
