package prettylog

import (
	"bytes"
	"log/slog"
	"sync"
)

// Handler represents a type that contains a slog.Handler, a function, a bytes.Buffer, and a sync.Mutex.
type Handler struct {
	h           slog.Handler
	r           func([]string, slog.Attr) slog.Attr
	b           *bytes.Buffer
	m           *sync.Mutex
	handlerType string
}

// NewHandler is a function that creates a new instance of the Handler struct
func NewHandler(handlerType string, opts *slog.HandlerOptions) *Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	b := &bytes.Buffer{}

	var h slog.Handler
	switch handlerType {
	case "text":
		h = slog.NewTextHandler(b, &slog.HandlerOptions{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: suppressDefaults(opts.ReplaceAttr),
		})
	case "json":
		h = slog.NewJSONHandler(b, &slog.HandlerOptions{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: suppressDefaults(opts.ReplaceAttr),
		})
	}

	return &Handler{
		handlerType: handlerType,
		h:           h,
		r:           opts.ReplaceAttr,
		b:           b,
		m:           &sync.Mutex{},
	}
}
