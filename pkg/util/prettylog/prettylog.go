package prettylog

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
)

type baseAttributes struct {
	ctx                   context.Context
	level, timestamp, msg string
	r                     slog.Record
}

// Enabled checks if the log handler is enabled for the given log level
func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

// WithAttrs returns a new log handler with additional attributes
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{h: h.h.WithAttrs(attrs), b: h.b, r: h.r, m: h.m}
}

// WithGroup returns a new log handler with the specified group name
func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{h: h.h.WithGroup(name), b: h.b, r: h.r, m: h.m}
}

// Handle is a method of the Handler struct that handles a slog.Record and returns an error
func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	baseAttributes := h.processBaseAttrs(ctx, r)
	switch h.handlerType {
	case "json":
		return h.handleJSON(baseAttributes)
	case "text":
		return h.handleText(baseAttributes)
	default:
		return fmt.Errorf("unknown handler type: %s", h.handlerType)
	}
}

func (h *Handler) processBaseAttrs(ctx context.Context, r slog.Record) baseAttributes {
	var level string
	levelAttr := slog.Attr{
		Key:   slog.LevelKey,
		Value: slog.AnyValue(r.Level),
	}
	if h.r != nil {
		levelAttr = h.r([]string{}, levelAttr)
	}

	if !levelAttr.Equal(slog.Attr{}) {
		level = levelAttr.Value.String()[:3]

		if r.Level <= slog.LevelDebug {
			level = colorize(lightGray, level)
		} else if r.Level <= slog.LevelInfo {
			level = colorize(cyan, level)
		} else if r.Level < slog.LevelWarn {
			level = colorize(lightBlue, level)
		} else if r.Level < slog.LevelError {
			level = colorize(lightYellow, level)
		} else if r.Level <= slog.LevelError+1 {
			level = colorize(lightRed, level)
		} else if r.Level > slog.LevelError+1 {
			level = colorize(lightMagenta, level)
		}
	}

	var timestamp string
	timeAttr := slog.Attr{
		Key:   slog.TimeKey,
		Value: slog.StringValue(r.Time.Format(timeFormat)),
	}
	if h.r != nil {
		timeAttr = h.r([]string{}, timeAttr)
	}
	if !timeAttr.Equal(slog.Attr{}) {
		timestamp = colorize(lightGray, timeAttr.Value.String())
	}

	var msg string
	msgAttr := slog.Attr{
		Key:   slog.MessageKey,
		Value: slog.StringValue(r.Message),
	}
	if h.r != nil {
		msgAttr = h.r([]string{}, msgAttr)
	}
	if !msgAttr.Equal(slog.Attr{}) {
		msg = colorize(white, msgAttr.Value.String())
	}

	return baseAttributes{
		ctx:       ctx,
		level:     level,
		timestamp: timestamp,
		msg:       msg,
		r:         r,
	}
}

func (h *Handler) handleJSON(baseAttrs baseAttributes) error {
	// JSON handler logic
	attrs, err := h.computeAttrs(baseAttrs.ctx, baseAttrs.r)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(attrs, "", "  ")
	if err != nil {
		return fmt.Errorf("error when marshaling attrs: %w", err)
	}

	out := strings.Builder{}
	if len(baseAttrs.timestamp) > 0 {
		out.WriteString(baseAttrs.timestamp)
		out.WriteString(" ")
	}
	if len(baseAttrs.level) > 0 {
		out.WriteString(baseAttrs.level)
		out.WriteString(" ")
	}
	if len(baseAttrs.msg) > 0 {
		out.WriteString(baseAttrs.msg)
		out.WriteString(" ")
	}
	if len(bytes) > 0 {
		out.WriteString(colorize(darkGray, string(bytes)))
	}

	fmt.Println(out.String())

	return nil
}

func (h *Handler) handleText(baseAttrs baseAttributes) error {
	// Text handler logic
	attrs := make(map[string]interface{}, baseAttrs.r.NumAttrs())
	baseAttrs.r.Attrs(func(a slog.Attr) bool {
		value := a.Value.Any()

		// Handle nil value
		if value == nil {
			return true // Skip this attribute
		}

		// Start the recursion with an empty prefix
		flattenAttributes(a.Key, value, attrs)

		return true
	})

	out := strings.Builder{}
	if len(baseAttrs.timestamp) > 0 {
		out.WriteString(baseAttrs.timestamp)
		out.WriteString(" ")
	}
	if len(baseAttrs.level) > 0 {
		out.WriteString(baseAttrs.level)
		out.WriteString(" ")
	}
	if len(baseAttrs.msg) > 0 {
		out.WriteString(baseAttrs.msg)
		out.WriteString(" ")
	}

	for key, value := range attrs {
		key = colorize(lightRed, key)
		out.WriteString(fmt.Sprintf("%s=%v ", key, value))
	}

	fmt.Println(out.String())

	return nil
}

func (h *Handler) computeAttrs(ctx context.Context, r slog.Record) (map[string]any, error) {
	h.m.Lock()
	defer func() {
		h.b.Reset()
		h.m.Unlock()
	}()
	if err := h.h.Handle(ctx, r); err != nil {
		return nil, fmt.Errorf("error when calling inner handler's Handle: %w", err)
	}

	var attrs map[string]any
	err := json.Unmarshal(h.b.Bytes(), &attrs)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshaling inner handler's Handle result: %w", err)
	}

	return attrs, nil
}
