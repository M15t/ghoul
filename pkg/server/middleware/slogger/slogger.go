package slogger

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	dblogger "github.com/M15t/ghoul/pkg/util/db/logger"
	"github.com/M15t/ghoul/pkg/util/threadsafe"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel/trace"
)

// custom
const (
	customAttributesCtxKey = "slog-echo.custom-attributes"
)

// custom var
var (
	timeFormat          = "Jan 02 15:04:05.000"
	RequestBodyMaxSize  = 64 * 1024 // 64KB
	ResponseBodyMaxSize = 64 * 1024 // 64KB

	HiddenRequestHeaders = map[string]struct{}{
		"authorization": {},
		"cookie":        {},
		"set-cookie":    {},
		"x-auth-token":  {},
		"x-csrf-token":  {},
		"x-xsrf-token":  {},
	}
	HiddenResponseHeaders = map[string]struct{}{
		"set-cookie": {},
	}

	sensitiveKeys = []string{"password", "key", "token", "username", "email", "phone", "mobile"}
)

// Config provide configurations
type Config struct {
	DefaultLevel     slog.Level
	ClientErrorLevel slog.Level
	ServerErrorLevel slog.Level

	WithUserAgent      bool
	WithRequestID      bool
	WithRequestBody    bool
	WithRequestHeader  bool
	WithResponseBody   bool
	WithResponseHeader bool
	WithSpanID         bool
	WithTraceID        bool
	WithDBQueries      bool

	Filters []Filter
}

// New returns a echo.MiddlewareFunc (middleware) that logs requests using slog.
//
// Requests with errors are logged using slog.Error().
// Requests without errors are logged using slog.Info().
func New(logger *slog.Logger) echo.MiddlewareFunc {
	return NewWithConfig(logger, Config{
		DefaultLevel:     slog.LevelInfo,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,

		WithUserAgent:      false,
		WithRequestID:      true,
		WithRequestBody:    false,
		WithRequestHeader:  false,
		WithResponseBody:   false,
		WithResponseHeader: false,
		WithSpanID:         false,
		WithTraceID:        false,
		WithDBQueries:      false,

		Filters: []Filter{},
	})
}

// NewWithFilters returns a echo.MiddlewareFunc (middleware) that logs requests using slog.
//
// Requests with errors are logged using slog.Error().
// Requests without errors are logged using slog.Info().
func NewWithFilters(logger *slog.Logger, filters ...Filter) echo.MiddlewareFunc {
	return NewWithConfig(logger, Config{
		DefaultLevel:     slog.LevelInfo,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,

		WithUserAgent:      false,
		WithRequestID:      true,
		WithRequestBody:    false,
		WithRequestHeader:  false,
		WithResponseBody:   false,
		WithResponseHeader: false,
		WithSpanID:         false,
		WithTraceID:        false,
		WithDBQueries:      false,

		Filters: filters,
	})
}

// NewWithConfig returns a echo.HandlerFunc (middleware) that logs requests using slog.
func NewWithConfig(logger *slog.Logger, config Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()
			start := time.Now()
			path := req.URL.Path
			query := req.URL.RawQuery

			ctx := req.Context()
			dbWriter := threadsafe.NewSimpleSlice([]string{})
			ctx = dblogger.WithContextGormLogger(ctx, dbWriter)
			// Replace the existing context in the request
			c.SetRequest(c.Request().WithContext(ctx))

			params := map[string]string{}
			for i, k := range c.ParamNames() {
				params[k] = c.ParamValues()[i]
			}

			// dump request body
			br := newBodyReader(req.Body, RequestBodyMaxSize, config.WithRequestBody)
			req.Body = br

			// dump response body
			bw := newBodyWriter(res.Writer, ResponseBodyMaxSize, config.WithResponseBody)
			res.Writer = bw

			err = next(c)

			if err != nil {
				c.Error(err)
			}

			status := res.Status
			method := req.Method
			end := time.Now()
			latency := end.Sub(start)
			userAgent := req.UserAgent()
			ip := c.RealIP()

			httpErr := new(echo.HTTPError)
			if err != nil && errors.As(err, &httpErr) {
				status = httpErr.Code
				if msg, ok := httpErr.Message.(string); ok {
					err = errors.New(msg)
				}
			}

			baseAttributes := []slog.Attr{
				slog.String("method", method),
				slog.String("path", path),
				slog.Int("status", status),
				slog.String("latency", latency.String()),
			}

			requestAttributes := []slog.Attr{
				slog.String("time", start.Format(timeFormat)),
				slog.String("ip", ip),
			}

			responseAttributes := []slog.Attr{
				slog.String("time", end.Format(timeFormat)),
			}

			if config.WithRequestID {
				requestID := req.Header.Get(echo.HeaderXRequestID)
				if requestID == "" {
					requestID = res.Header().Get(echo.HeaderXRequestID)
				}
				if requestID != "" {
					baseAttributes = append(baseAttributes, slog.String("id", requestID))
				}
			}

			if query != "" {
				requestAttributes = append(requestAttributes, slog.String("query", query))
			}

			if len(params) > 0 {
				requestAttributes = append(requestAttributes, slog.Any("params", params))
			}

			// otel
			if config.WithTraceID {
				traceID := trace.SpanFromContext(c.Request().Context()).SpanContext().TraceID().String()
				baseAttributes = append(baseAttributes, slog.String("trace-id", traceID))
			}
			if config.WithSpanID {
				spanID := trace.SpanFromContext(c.Request().Context()).SpanContext().SpanID().String()
				baseAttributes = append(baseAttributes, slog.String("span-id", spanID))
			}

			// db queries
			if config.WithDBQueries {
				if len(dbWriter.All()) > 0 {
					baseAttributes = append(baseAttributes, slog.Any("queries", dbWriter.All()))
				}
			}

			// request body
			if config.WithRequestBody {
				// proceed body dump
				reqBody := br.body.Bytes()
				if strings.Contains(c.Request().Header.Get("Content-Type"), "application/json") && len(reqBody) > 0 {
					requestAttributes = append(requestAttributes, slog.Any("body", prettyJSON(reqBody)))
				}
			}

			// request headers
			if config.WithRequestHeader {
				kv := []any{}

				for k, v := range c.Request().Header {
					if _, found := HiddenRequestHeaders[strings.ToLower(k)]; found {
						continue
					}
					kv = append(kv, slog.Any(k, v))
				}

				requestAttributes = append(requestAttributes, slog.Group("header", kv...))
			}

			if config.WithUserAgent {
				requestAttributes = append(requestAttributes, slog.String("user-agent", userAgent))
			}

			xForwardedFor, ok := c.Get(echo.HeaderXForwardedFor).(string)
			if ok && len(xForwardedFor) > 0 {
				ips := lo.Map(strings.Split(xForwardedFor, ","), func(ip string, _ int) string {
					return strings.TrimSpace(ip)
				})
				requestAttributes = append(requestAttributes, slog.Any("x-forwarded-for", ips))
			}

			// response body body
			if config.WithResponseBody {
				// proceed body dump
				resBody := bw.body.Bytes()

				if strings.Contains(c.Response().Header().Get("Content-Type"), "application/json") && len(resBody) > 0 {
					responseAttributes = append(responseAttributes, slog.Any("body", prettyJSON(resBody)))
				}
			}

			// response headers
			if config.WithResponseHeader {
				kv := []any{}

				for k, v := range c.Response().Header() {
					if _, found := HiddenResponseHeaders[strings.ToLower(k)]; found {
						continue
					}
					kv = append(kv, slog.Any(k, v))
				}

				responseAttributes = append(responseAttributes, slog.Group("header", kv...))
			}

			attributes := baseAttributes

			if config.WithRequestBody {
				attributes = append(attributes, slog.Attr{
					Key:   "request",
					Value: slog.GroupValue(requestAttributes...),
				})
			}

			if config.WithResponseBody {
				attributes = append(attributes, slog.Attr{
					Key:   "response",
					Value: slog.GroupValue(responseAttributes...),
				})
			}

			// custom context values
			if v := c.Get(customAttributesCtxKey); v != nil {
				switch attrs := v.(type) {
				case []slog.Attr:
					attributes = append(attributes, attrs...)
				}
			}

			for _, filter := range config.Filters {
				if !filter(c) {
					return
				}
			}

			level := config.DefaultLevel
			msg := "OK"
			if status >= http.StatusInternalServerError {
				level = config.ServerErrorLevel
				if err != nil {
					msg = err.Error()
				} else {
					msg = http.StatusText(status)
				}
			} else if status >= http.StatusBadRequest && status < http.StatusInternalServerError {
				level = config.ClientErrorLevel
				msg = http.StatusText(status)
			}

			logger.LogAttrs(c.Request().Context(), level, msg, attributes...)

			return
		}
	}
}

// AddCustomAttributes adds custom slog.Attr
func AddCustomAttributes(c echo.Context, attr slog.Attr) {
	v := c.Get(customAttributesCtxKey)
	if v == nil {
		c.Set(customAttributesCtxKey, []slog.Attr{attr})
		return
	}

	switch attrs := v.(type) {
	case []slog.Attr:
		c.Set(customAttributesCtxKey, append(attrs, attr))
	}
}
