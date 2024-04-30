package slogger

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MiddlewareFunc is returned
func TestMiddlewareFuncIsReturned(t *testing.T) {
	// Initialize logger and config
	logger := &slog.Logger{}
	config := Config{
		DefaultLevel:       slog.LevelInfo,
		ClientErrorLevel:   slog.LevelWarn,
		ServerErrorLevel:   slog.LevelError,
		WithUserAgent:      false,
		WithRequestID:      true,
		WithRequestBody:    false,
		WithRequestHeader:  false,
		WithResponseBody:   false,
		WithResponseHeader: false,
		WithSpanID:         false,
		WithTraceID:        false,
		WithDBQueries:      false,
		Filters:            []Filter{},
	}

	// Call the function under test
	middlewareFunc := NewWithConfig(logger, config)

	// Assert that the returned value is a MiddlewareFunc
	assert.NotNil(t, middlewareFunc)
}
