package echoapp_middlewares

import (
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/random"
	"github.com/sirupsen/logrus"
	"time"
)

type LoggingMiddlewareConfig struct {
	Name      string
	Skipper   middleware.Skipper
	Generator func() string
	// logger
	Logger *logrus.Entry
}

var (
	// DefaultRequestIDConfig is the default RequestID middleware config.
	DefaultLoggingMiddlewareConfig = LoggingMiddlewareConfig{
		Name:      "echo-app",
		Skipper:   middleware.DefaultSkipper,
		Generator: generatorId,
		Logger:    echoapp_util.NewDefaultEntry(),
	}
)

func NewLoggingMiddleware(config LoggingMiddlewareConfig) echo.MiddlewareFunc {
	skipper := config.Skipper
	if config.Name == "" {
		config.Name = DefaultLoggingMiddlewareConfig.Name
	}
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultLoggingMiddlewareConfig.Skipper
	}
	if config.Generator == nil {
		config.Generator = generatorId
	}
	if config.Logger == nil {
		config.Logger = DefaultLoggingMiddlewareConfig.Logger
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if skipper(c) {
				return next(c)
			}

			start := time.Now()
			req := c.Request()
			rid := req.Header.Get(echo.HeaderXRequestID)
			if rid == "" {
				rid = config.Generator()
				req.Header.Set(echo.HeaderXRequestID, rid)
			}

			echoapp_util.AddRequestId(c, rid)
			echoapp_util.ToContext(c, config.Logger)

			fields := logrus.Fields{
				"Host":      req.Host,
				"Remote":    c.RealIP(),
				"Method":    req.Method,
				"URI":       req.RequestURI,
				"Referer":   req.Referer(),
				"UserAgent": req.UserAgent(),
			}
			logger := echoapp_util.ExtractEntry(c).
				WithField("App", config.Name).
				WithFields(fields)
			echoapp_util.ToContext(c, logger)
			err := next(c)
			// in case any step changed the logger context
			logger = echoapp_util.ExtractEntry(c)
			latency := time.Since(start)
			resp := c.Response()
			if err != nil {
				logger.WithField("status", resp.Status).WithField("latency", latency.Nanoseconds()/int64(time.Millisecond)).
					Error(err.Error())
			} else {
				logger.WithField("status", resp.Status).WithField("latency", latency.Nanoseconds()/int64(time.Millisecond)).
					Info("LoggingMiddleware request over.")
			}

			return nil
		}
	}
}

func generatorId() string {
	return random.String(32)
}
