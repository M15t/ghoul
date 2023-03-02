package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/M15t/ghoul/pkg/server/middleware/secure"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
)

// Config represents server specific config
type Config struct {
	Stage        string
	Port         int
	ReadTimeout  int
	WriteTimeout int
	Debug        bool
	AllowOrigins []string
}

var (
	// DefaultConfig for the API server
	DefaultConfig = Config{
		Stage:        "development",
		Port:         8080,
		ReadTimeout:  10,
		WriteTimeout: 5,
		Debug:        true,
		AllowOrigins: []string{"*"},
	}

	version   = "dev" // sha1 revision used to build the server
	buildTime = "now" // when the server was built
)

func (c *Config) fillDefaults() {
	if c.Stage == "" {
		c.Stage = DefaultConfig.Stage
	}
	if c.Port == 0 {
		c.Port = DefaultConfig.Port
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = DefaultConfig.ReadTimeout
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = DefaultConfig.WriteTimeout
	}
	if c.AllowOrigins == nil && len(c.AllowOrigins) == 0 {
		c.AllowOrigins = DefaultConfig.AllowOrigins
	}
}

var echoLambda *echoadapter.EchoLambda

// New instantates new Echo server
func New(cfg *Config) *echo.Echo {
	cfg.fillDefaults()
	e := echo.New()
	e.Use(middleware.Logger(), middleware.Recover(), secure.Headers(), secure.CORS(&secure.Config{AllowOrigins: cfg.AllowOrigins}))
	e.GET("/", healthCheck)
	e.Validator = NewValidator()
	e.HTTPErrorHandler = NewErrorHandler(e).Handle
	e.Binder = NewBinder()
	e.Debug = cfg.Debug
	if e.Debug {
		e.Logger.SetLevel(log.DEBUG)
		e.Use(secure.BodyDump())
	} else {
		e.Logger.SetLevel(log.ERROR)
	}
	e.Server.Addr = fmt.Sprintf(":%d", cfg.Port)
	e.Server.ReadTimeout = time.Duration(cfg.ReadTimeout) * time.Minute
	e.Server.WriteTimeout = time.Duration(cfg.WriteTimeout) * time.Minute

	return e
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":     "ok",
		"version":    version,
		"build_time": buildTime,
	})
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return echoLambda.ProxyWithContext(ctx, req)
}

// Start starts echo server
func Start(e *echo.Echo, isDevelopment bool) {
	// graceful shutdown for dev environment
	if isDevelopment {
		// Start server
		go func() {
			if err := e.StartServer(e.Server); err != nil {
				if err == http.ErrServerClosed {
					e.Logger.Info("shutting down the server")
				} else {
					e.Logger.Errorf("error shutting down the server: ", err)
				}
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 10 seconds.
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}
	} else {
		// Hide verbose logs and start server normally
		e.HideBanner = true
		e.HidePort = true
		e.Logger.Fatal(e.StartServer(e.Server))

		// User echo adapter for Lambda
		echoLambda = echoadapter.New(e)
		lambda.Start(handler)
	}
}
