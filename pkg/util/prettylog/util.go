package prettylog

import (
	"fmt"
	"log/slog"
	"strconv"
)

const (
	timeFormat = "Jan 02 15:04:05.000"

	reset = "\033[0m"

	black        = 30
	red          = 31
	green        = 32
	yellow       = 33
	blue         = 34
	magenta      = 35
	cyan         = 36
	lightGray    = 37
	darkGray     = 90
	lightRed     = 91
	lightGreen   = 92
	lightYellow  = 93
	lightBlue    = 94
	lightMagenta = 95
	lightCyan    = 96
	white        = 97
)

func colorize(colorCode int, v string) string {
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), v, reset)
}

func suppressDefaults(
	next func([]string, slog.Attr) slog.Attr,
) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey ||
			a.Key == slog.LevelKey ||
			a.Key == slog.MessageKey {
			return slog.Attr{}
		}
		if next == nil {
			return a
		}
		return next(groups, a)
	}
}

func flattenAttributes(prefix string, value interface{}, attrs map[string]interface{}) {
	switch v := value.(type) {
	case []slog.Attr:
		// Flatten nested attributes
		for _, c := range v {
			// Construct the key as "prefix.child"
			newKey := fmt.Sprintf("%s.%s", prefix, c.Key)
			flattenAttributes(newKey, c.Value, attrs)
		}
	default:
		// If value is not a slice, store it directly
		attrs[prefix] = value
	}
}
