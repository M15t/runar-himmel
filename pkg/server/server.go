package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"runar-himmel/pkg/server/middleware/secure"

	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
)

// Adapter for Echo when running on AWS Lambda
var echoLambda *echoadapter.EchoLambdaV2

// Config represents server specific config
type Config struct {
	Debug bool

	// The port for the http server to listen on
	Port int
	// ReadHeaderTimeout is the amount of time allowed to read request headers.
	ReadHeaderTimeout int
	// ReadTimeout is the maximum duration for reading the entire request, including the body
	ReadTimeout int
	// WriteTimeout is the maximum duration before timing out writes of the response
	WriteTimeout int

	// CORS settings
	AllowOrigins []string
	// Maximum allowed size for a request body. ex: 512KB
	BodyLimit string
	// To skip the BodyLimit check. If not provided, skip when path contains `/upload` or has `/admin` prefix
	BodyLimitSkipper middleware.Skipper
	// The `Content-Security-Policy` header providing security against XSS and other code injection attacks.
	// Sample for production: `default-src 'self'`
	ContentSecurityPolicy string
}

var (
	// DefaultConfig for the API server
	DefaultConfig = Config{
		Port:              8080,
		ReadHeaderTimeout: 10,
		ReadTimeout:       30,
		WriteTimeout:      60,
		Debug:             true,
		AllowOrigins:      []string{"*"},
		BodyLimit:         "512KB",
		BodyLimitSkipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "/upload") || strings.HasPrefix(c.Path(), "/admin")
		},
	}
)

func (c *Config) fillDefaults() {
	if c.Port == 0 {
		c.Port = DefaultConfig.Port
	}
	if c.ReadHeaderTimeout == 0 {
		c.ReadHeaderTimeout = DefaultConfig.ReadHeaderTimeout
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = DefaultConfig.ReadTimeout
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = DefaultConfig.WriteTimeout
	}
	if len(c.AllowOrigins) == 0 {
		c.AllowOrigins = DefaultConfig.AllowOrigins
	}
	if c.BodyLimit == "" {
		c.BodyLimit = DefaultConfig.BodyLimit
	}
	if c.BodyLimitSkipper == nil {
		c.BodyLimitSkipper = DefaultConfig.BodyLimitSkipper
	}
}

// New instantates new Echo server
func New(cfg *Config) *echo.Echo {
	cfg.fillDefaults()
	e := echo.New()
	e.Validator = NewValidator()
	e.HTTPErrorHandler = NewErrorHandler(e).Handle
	e.Binder = NewBinder()
	e.Debug = cfg.Debug
	e.Server.Addr = fmt.Sprintf(":%d", cfg.Port)
	e.Server.ReadHeaderTimeout = time.Duration(cfg.ReadHeaderTimeout) * time.Second
	e.Server.ReadTimeout = time.Duration(cfg.ReadTimeout) * time.Second
	e.Server.WriteTimeout = time.Duration(cfg.WriteTimeout) * time.Second

	e.Use(middleware.Recover(), middleware.RequestID(), middleware.Logger())
	e.Use(middleware.BodyLimitWithConfig(middleware.BodyLimitConfig{
		Limit:   cfg.BodyLimit,
		Skipper: cfg.BodyLimitSkipper,
	}))

	if e.Debug {
		e.Logger.SetLevel(log.DEBUG)
		e.Use(secure.BodyDump())
	} else {
		e.Logger.SetLevel(log.INFO)
	}

	// add security & cors middlewares
	e.Use(secure.Headers(cfg.ContentSecurityPolicy), secure.SimpleCORS(cfg.AllowOrigins))

	return e
}

// Start starts echo server with graceful shutdown process
func Start(e *echo.Echo, isLambda bool) {
	// hide verbose logs
	e.HideBanner = true

	if !isLambda {
		go func() {
			if err := e.StartServer(e.Server); err != nil {
				if err == http.ErrServerClosed {
					fmt.Printf("⇨ http server stopped\n")
				} else {
					fmt.Printf("⇨ http server starting error: %v\n", err)
				}
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
		// Use a buffered channel to avoid missing signals as recommended for signal.Notify
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit

		// received signal, shutting down...
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			fmt.Printf("⇨ http server shutting down error: %v\n", err)
		}
	} else {
		e.HidePort = true

		echoLambda = echoadapter.NewV2(e)
		lambda.Start(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
			return echoLambda.ProxyWithContext(ctx, req)
		})
	}
}
