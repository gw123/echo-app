package echoapp_middlewares

import (
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/random"
	"github.com/sirupsen/logrus"
)

type (
	// RequestIDConfig defines the config for RequestID middleware.
	ContextConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper
		// Generator defines a function to generate an ID.
		// Optional. Default value random.String(32).
		Generator func() string
		// logger
		Logger *logrus.Entry
	}
)

var (
	// DefaultRequestIDConfig is the default RequestID middleware config.
	DefaultRequestIDConfig = ContextConfig{
		Skipper:   middleware.DefaultSkipper,
		Generator: generator,
		Logger:    echoapp_util.NewDefaultEntry(),
	}
)

// RequestID returns a X-Request-ID middleware.
func RequestID() echo.MiddlewareFunc {
	return RequestIDWithConfig(DefaultRequestIDConfig)
}

// RequestIDWithConfig returns a X-Request-ID middleware with config.
func RequestIDWithConfig(config ContextConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRequestIDConfig.Skipper
	}
	if config.Generator == nil {
		config.Generator = generator
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()

			rid := req.Header.Get(echo.HeaderXRequestID)
			if rid == "" {
				rid = config.Generator()
				req.Header.Set(echo.HeaderXRequestID, rid)
			}

			echoapp_util.AddRequestId(c, rid)
			echoapp_util.ToContext(c, config.Logger)
			return next(c)
		}
	}
}

func generator() string {
	return random.String(32)
}
