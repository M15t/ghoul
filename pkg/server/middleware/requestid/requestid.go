package requestid

import (
	"context"

	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/labstack/echo/v4"
)

// Define a custom type for the request ID key
type requestIDKey string

// RequestIDKey define a constant string for the request ID key
const RequestIDKey requestIDKey = "custom.requestID"

// Config defines the config for RequestID middleware.
type Config struct {
	// Generator defines a function to generate an ID.
	// Optional. Defaults to generator for random string of length 32.
	Generator func() string

	// RequestIDHandler defines a function which is executed for a request id.
	RequestIDHandler func(echo.Context, string)

	// TargetHeader defines what header to look for to populate the id
	TargetHeader string
}

// DefaultRequestIDConfig is the default RequestID middleware config.
var DefaultRequestIDConfig = Config{
	Generator:    generator,
	TargetHeader: echo.HeaderXRequestID,
}

// New returns a X-Request-ID middleware.
func New() echo.MiddlewareFunc {
	return NewWithConfig(DefaultRequestIDConfig)
}

// NewWithConfig returns a X-Request-ID middleware with config.
func NewWithConfig(config Config) echo.MiddlewareFunc {
	// Defaults
	if config.Generator == nil {
		config.Generator = generator
	}
	if config.TargetHeader == "" {
		config.TargetHeader = echo.HeaderXRequestID
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			var rid string

			// Check if AWS Lambda context is available and if AwsRequestID is not empty
			if lambdaCtx, ok := core.GetRuntimeContextFromContextV2(req.Context()); ok && lambdaCtx.AwsRequestID != "" {
				rid = lambdaCtx.AwsRequestID
			} else {
				// Check the value of TargetHeader in the request headers
				rid = req.Header.Get(config.TargetHeader)
			}

			// Generate a random request ID if both AWS Lambda context and TargetHeader value are empty
			if rid == "" {
				rid = config.Generator()
			}
			res.Header().Set(config.TargetHeader, rid)
			if config.RequestIDHandler != nil {
				config.RequestIDHandler(c, rid)
			}

			// Create a new context with the request ID
			ctx := WithContextRequestID(req.Context(), rid)
			// Replace the existing context in the request
			c.SetRequest(req.WithContext(ctx))

			return next(c)
		}
	}
}

// WithContextRequestID takes a context and request id as inputs and returns a new context with a value
func WithContextRequestID(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, RequestIDKey, rid)
}

// GetContextRequestID takes a context as input and returns a pointer to a string of strings
func GetContextRequestID(ctx context.Context) string {
	if value, ok := ctx.Value(RequestIDKey).(string); ok {
		return value
	}

	return ""
}

func generator() string {
	return randomString(32)
}
