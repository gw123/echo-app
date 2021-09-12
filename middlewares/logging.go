package echoapp_middlewares

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"

	"github.com/gw123/glog"
	glogCommon "github.com/gw123/glog/common"

	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/random"
	"github.com/sirupsen/logrus"
)

type LoggingMiddlewareConfig struct {
	Name               string
	Skipper            middleware.Skipper
	Generator          func() string
	Logger             glogCommon.Logger
	EnableTrace        bool
	TraceAgentHostPort string
	CloseChannel       chan struct{}
}

var (
	// DefaultRequestIDConfig is the default RequestID middleware config.
	DefaultLoggingMiddlewareConfig = LoggingMiddlewareConfig{
		Name:               "echo-app",
		Skipper:            middleware.DefaultSkipper,
		Generator:          generatorId,
		Logger:             glog.DefaultLogger(),
		EnableTrace:        true,
		TraceAgentHostPort: "127.0.0.1:6831",
		CloseChannel:       make(chan struct{}),
	}
)

func initJaeger(middlewareConfig LoggingMiddlewareConfig) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: middlewareConfig.TraceAgentHostPort,
		},
	}
	tracer, closer, err := cfg.New(middlewareConfig.Name, config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("Error: connot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

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

	var tracer opentracing.Tracer
	var closer io.Closer
	if config.EnableTrace {
		// todo 实现优雅关闭
		tracer, closer = initJaeger(config)
		if config.CloseChannel != nil {
			go func() {
				<-config.CloseChannel
				closer.Close()
			}()
		}
		opentracing.SetGlobalTracer(tracer)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if skipper(c) {
				return next(c)
			}

			start := time.Now()
			req := c.Request()
			if strings.Contains(req.RequestURI, ".js") ||
				strings.Contains(req.RequestURI, ".css") {
				return next(c)
			}
			ctx := context.Background()
			echoapp_util.SetCtxContext(c, ctx)

			if config.EnableTrace {
				span := tracer.StartSpan("span_root")
				ctx = opentracing.ContextWithSpan(ctx, span)
				defer span.Finish()
			}

			rid := req.Header.Get(echo.HeaderXRequestID)
			if rid == "" {
				rid = config.Generator()
				req.Header.Set(echo.HeaderXRequestID, rid)
			}
			//日志太多了暂时屏蔽掉 requestid
			//echoapp_util.AddRequestId(c, rid)
			//echoapp_util.ToContext(c, config.Logger)

			fields := logrus.Fields{
				//"host":   req.Host,
				"remote": c.RealIP(),
				//"method": req.Method,
				//"Referer": req.Referer(),
				//"UserAgent": req.UserAgent(),
			}

			logger := config.Logger.
				WithFields(fields).
				WithField(glogCommon.KeyTraceID, rid).
				WithField(glogCommon.KeyPathname, c.Request().URL.Path)

			echoapp_util.ToContext(c, logger)
			err := next(c)
			// in case any step changed the logger context
			logger = echoapp_util.ExtractEntry(c)
			latency := time.Since(start)
			resp := c.Response()
			if err != nil {
				logger.WithField("status", resp.Status).WithField("latency", latency.Nanoseconds()/int64(time.Millisecond)).
					WithField("error", err.Error()).Error("log middleware err")
			} else {
				logger.WithField("status", resp.Status).WithField("latency", latency.Nanoseconds()/int64(time.Millisecond)).
					Info("log middleware success")
			}

			return nil
		}
	}
}

func generatorId() string {
	return random.String(32)
}
