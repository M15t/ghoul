package slogger

import (
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
)

// Filter is core filter
type Filter func(ctx echo.Context) bool

// AcceptMethod filters allow by method
func AcceptMethod(methods ...string) Filter {
	return func(c echo.Context) bool {
		reqMethod := strings.ToLower(c.Request().Method)

		for _, method := range methods {
			if strings.ToLower(method) == reqMethod {
				return true
			}
		}

		return false
	}
}

// IgnoreMethod filters disallow by method
func IgnoreMethod(methods ...string) Filter {
	return func(c echo.Context) bool {
		reqMethod := strings.ToLower(c.Request().Method)

		for _, method := range methods {
			if strings.ToLower(method) == reqMethod {
				return false
			}
		}

		return true
	}
}

// AcceptStatus filters allow by status
func AcceptStatus(statuses ...int) Filter {
	return func(c echo.Context) bool {
		for _, status := range statuses {
			if status == c.Response().Status {
				return true
			}
		}

		return false
	}
}

// IgnoreStatus filters disallow by status
func IgnoreStatus(statuses ...int) Filter {
	return func(c echo.Context) bool {
		for _, status := range statuses {
			if status == c.Response().Status {
				return false
			}
		}

		return true
	}
}

// AcceptStatusGreaterThan filters allow by status gt
func AcceptStatusGreaterThan(status int) Filter {
	return func(c echo.Context) bool {
		return c.Response().Status > status
	}
}

// IgnoreStatusLessThan filters disallow by status lt
func IgnoreStatusLessThan(status int) Filter {
	return func(c echo.Context) bool {
		return c.Response().Status < status
	}
}

// AcceptStatusGreaterThanOrEqual filters allow by status gte
func AcceptStatusGreaterThanOrEqual(status int) Filter {
	return func(c echo.Context) bool {
		return c.Response().Status >= status
	}
}

// IgnoreStatusLessThanOrEqual filters disallow by status lte
func IgnoreStatusLessThanOrEqual(status int) Filter {
	return func(c echo.Context) bool {
		return c.Response().Status <= status
	}
}

// AcceptPath filters allow by path
func AcceptPath(urls ...string) Filter {
	return func(c echo.Context) bool {
		for _, url := range urls {
			if c.Request().URL.Path == url {
				return true
			}
		}

		return false
	}
}

// IgnorePath filters disallow by path
func IgnorePath(urls ...string) Filter {
	return func(c echo.Context) bool {
		for _, url := range urls {
			if c.Request().URL.Path == url {
				return false
			}
		}

		return true
	}
}

// AcceptPathContains filters allow by path contains
func AcceptPathContains(parts ...string) Filter {
	return func(c echo.Context) bool {
		for _, part := range parts {
			if strings.Contains(c.Request().URL.Path, part) {
				return true
			}
		}

		return false
	}
}

// IgnorePathContains filters disallow by path contains
func IgnorePathContains(parts ...string) Filter {
	return func(c echo.Context) bool {
		for _, part := range parts {
			if strings.Contains(c.Request().URL.Path, part) {
				return false
			}
		}

		return true
	}
}

// AcceptPathPrefix filters allow by path prefix
func AcceptPathPrefix(prefixs ...string) Filter {
	return func(c echo.Context) bool {
		for _, prefix := range prefixs {
			if strings.HasPrefix(c.Request().URL.Path, prefix) {
				return true
			}
		}

		return false
	}
}

// IgnorePathPrefix filters disallow by path prefix
func IgnorePathPrefix(prefixs ...string) Filter {
	return func(c echo.Context) bool {
		for _, prefix := range prefixs {
			if strings.HasPrefix(c.Request().URL.Path, prefix) {
				return false
			}
		}

		return true
	}
}

// AcceptPathSuffix filters allow by path suffix
func AcceptPathSuffix(prefixs ...string) Filter {
	return func(c echo.Context) bool {
		for _, prefix := range prefixs {
			if strings.HasPrefix(c.Request().URL.Path, prefix) {
				return true
			}
		}

		return false
	}
}

// IgnorePathSuffix filters disallow by path suffix
func IgnorePathSuffix(suffixs ...string) Filter {
	return func(c echo.Context) bool {
		for _, suffix := range suffixs {
			if strings.HasSuffix(c.Request().URL.Path, suffix) {
				return false
			}
		}

		return true
	}
}

// AcceptPathMatch filters allow by path matched regexp
func AcceptPathMatch(regs ...regexp.Regexp) Filter {
	return func(c echo.Context) bool {
		for _, reg := range regs {
			if reg.Match([]byte(c.Request().URL.Path)) {
				return true
			}
		}

		return false
	}
}

// IgnorePathMatch filters disallow by path matched regexp
func IgnorePathMatch(regs ...regexp.Regexp) Filter {
	return func(c echo.Context) bool {
		for _, reg := range regs {
			if reg.Match([]byte(c.Request().URL.Path)) {
				return false
			}
		}

		return true
	}
}

// AcceptHost filters allow by host
func AcceptHost(hosts ...string) Filter {
	return func(c echo.Context) bool {
		for _, host := range hosts {
			if c.Request().URL.Host == host {
				return true
			}
		}

		return false
	}
}

// IgnoreHost filters disallow by host
func IgnoreHost(hosts ...string) Filter {
	return func(c echo.Context) bool {
		for _, host := range hosts {
			if c.Request().URL.Host == host {
				return false
			}
		}

		return true
	}
}

// AcceptHostContains filters allow host contains
func AcceptHostContains(parts ...string) Filter {
	return func(c echo.Context) bool {
		for _, part := range parts {
			if strings.Contains(c.Request().URL.Host, part) {
				return true
			}
		}

		return false
	}
}

// IgnoreHostContains filters disallow host contains
func IgnoreHostContains(parts ...string) Filter {
	return func(c echo.Context) bool {
		for _, part := range parts {
			if strings.Contains(c.Request().URL.Host, part) {
				return false
			}
		}

		return true
	}
}

// AcceptHostPrefix filters allow host by prefix
func AcceptHostPrefix(prefixs ...string) Filter {
	return func(c echo.Context) bool {
		for _, prefix := range prefixs {
			if strings.HasPrefix(c.Request().URL.Host, prefix) {
				return true
			}
		}

		return false
	}
}

// IgnoreHostPrefix filters disallow host by prefix
func IgnoreHostPrefix(prefixs ...string) Filter {
	return func(c echo.Context) bool {
		for _, prefix := range prefixs {
			if strings.HasPrefix(c.Request().URL.Host, prefix) {
				return false
			}
		}

		return true
	}
}

// AcceptHostSuffix filters allow by host suffix
func AcceptHostSuffix(prefixs ...string) Filter {
	return func(c echo.Context) bool {
		for _, prefix := range prefixs {
			if strings.HasPrefix(c.Request().URL.Host, prefix) {
				return true
			}
		}

		return false
	}
}

// IgnoreHostSuffix filters disallow by host suffix
func IgnoreHostSuffix(suffixs ...string) Filter {
	return func(c echo.Context) bool {
		for _, suffix := range suffixs {
			if strings.HasSuffix(c.Request().URL.Host, suffix) {
				return false
			}
		}

		return true
	}
}

// AcceptHostMatch filter allow by host matched regex
func AcceptHostMatch(regs ...regexp.Regexp) Filter {
	return func(c echo.Context) bool {
		for _, reg := range regs {
			if reg.Match([]byte(c.Request().URL.Host)) {
				return true
			}
		}

		return false
	}
}

// IgnoreHostMatch filter disallow by host matched regex
func IgnoreHostMatch(regs ...regexp.Regexp) Filter {
	return func(c echo.Context) bool {
		for _, reg := range regs {
			if reg.Match([]byte(c.Request().URL.Host)) {
				return false
			}
		}

		return true
	}
}
