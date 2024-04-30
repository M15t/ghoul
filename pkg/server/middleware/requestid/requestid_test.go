package requestid

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// Creates a middleware function that returns a handler function.
func TestNewWithConfigMiddlewareFunc(t *testing.T) {
	// Initialize the config
	config := Config{}

	// Call the function under test
	middlewareFunc := NewWithConfig(config)

	// Assert that the returned value is a function
	assert.NotNil(t, middlewareFunc)

	// Create a mock echo context
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the middleware function
	handlerFunc := middlewareFunc(func(c echo.Context) error { return nil })

	// Assert that the handler function does not return an error
	assert.NoError(t, handlerFunc(c))

	// Assert that the response header contains the request ID
	assert.NotEmpty(t, rec.Header().Get(echo.HeaderXRequestID))
}

// The request ID handler function modifies the response header.
func TestRequestIDHandlerModifiesResponseHeader(t *testing.T) {
	// Initialize the config
	config := Config{}

	// Call the function under test
	middlewareFunc := NewWithConfig(config)

	// Assert that the returned value is a function
	assert.NotNil(t, middlewareFunc)

	// Create a mock echo context
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the middleware function
	handlerFunc := middlewareFunc(func(c echo.Context) error { return nil })

	// Assert that the handler function does not return an error
	assert.NoError(t, handlerFunc(c))

	// Assert that the response header contains the request ID
	assert.NotEmpty(t, rec.Header().Get(echo.HeaderXRequestID))
}
