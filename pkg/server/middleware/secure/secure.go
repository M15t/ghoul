package secure

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/thoas/go-funk"
)

// Config represents secure specific config
type Config struct {
	AllowOrigins []string
}

// Headers adds general security headers for basic security measures
func Headers() echo.MiddlewareFunc {
	return middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSMaxAge:            31536000,
		HSTSExcludeSubdomains: true,
		// ContentSecurityPolicy: "default-src 'self'",
	})
}

// CORS adds Cross-Origin Resource Sharing support
func CORS(cfg *Config) echo.MiddlewareFunc {
	allowOrigins := []string{"*"}
	if cfg != nil && cfg.AllowOrigins != nil {
		allowOrigins = cfg.AllowOrigins
	}

	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "PATCH", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           86400,
	})
}

// BodyDump prints out the request body for debugging purpose
func BodyDump() echo.MiddlewareFunc {
	return middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		secretFields := []string{"new_password", "old_password", "password", "access_token", "refresh_token"}
		contentType := c.Request().Header.Get("Content-Type")

		if len(reqBody) > 0 && contentType == "application/json" {
			var bodymap map[string]interface{}
			if err := json.Unmarshal(reqBody, &bodymap); err == nil {
				for i := 0; i < len(secretFields); i++ {
					if _, ok := bodymap[secretFields[i]]; ok {
						bodymap[secretFields[i]] = "********"
					}
				}
				reqBody, _ = json.Marshal(bodymap)
			}
			if funk.ContainsString([]string{"Content-Disposition: form-data"}, string(reqBody)) {
				c.Logger().Debug("Request Body: ", "multipart/form-data")
			}
			c.Logger().Debug("Request Body: ", string(reqBody))
		}

		if (c.Request().Method == "PATCH" || c.Request().Method == "POST") && len(resBody) > 0 {
			var bodymap map[string]interface{}
			if err := json.Unmarshal(resBody, &bodymap); err == nil {
				for i := 0; i < len(secretFields); i++ {
					if _, ok := bodymap[secretFields[i]]; ok {
						bodymap[secretFields[i]] = "********"
					}
				}
				resBody, _ = json.Marshal(bodymap)
			}
			if funk.ContainsString([]string{"Content-Disposition: form-data"}, string(reqBody)) {
				c.Logger().Debug("Request Body: ", "multipart/form-data")
			}
			c.Logger().Debug("Response Body: ", string(resBody))
		}
	})
}
