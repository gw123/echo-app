package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/app"
	"github.com/gw123/echo-app/controllers"
	echoapp_middlewares "github.com/gw123/echo-app/middlewares"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/cobra"
)

func startOrderServer() {
	echoapp_util.DefaultLogger().Info("开启order服务")
	//echoapp_util.DefaultLogger().Infof("%+v", echoapp.ConfigOpts)
	e := echo.New()
	e.HTTPErrorHandler = func(err error, ctx echo.Context) {
		ctx.JSON(http.StatusInternalServerError, map[string]string{"msg": err.Error()})
	}

	if echoapp.ConfigOpts.Asset.PublicRoot != "" {
		e.Static("/", echoapp.ConfigOpts.Asset.PublicRoot)
	}
	assetConfig := echoapp.ConfigOpts.Asset
	e.Renderer = echoapp_util.NewTemplateRenderer(assetConfig.ViewRoot, assetConfig.PublicHost, assetConfig.Version)
	origins := echoapp.ConfigOpts.OrderServer.Origins
	corsMiddleware := middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: origins,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType,
			echo.HeaderAccept, "x-requested-with", "authorization", "x-csrf-token", "Access-Control-Allow-Credentials"},
	})

	loggerMiddleware := echoapp_middlewares.NewLoggingMiddleware(echoapp_middlewares.LoggingMiddlewareConfig{
		Skipper: func(ctx echo.Context) bool {
			req := ctx.Request()
			return (req.RequestURI == "/" && req.Method == "HEAD") || (req.RequestURI == "/favicon.ico" && req.Method == "GET")
		},
	})
	e.Use(loggerMiddleware)
	//e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
	//	StackSize: 1 << 10, // 1 KB
	//}))

	//Actions
	orderSvr := app.MustGetOrderService()
	companySvr := app.MustGetCompanyService()
	limitMiddleware := echoapp_middlewares.NewLimitMiddlewares(middleware.DefaultSkipper, 100, 200)
	companyMiddleware := echoapp_middlewares.NewCompanyMiddlewares(middleware.DefaultSkipper, companySvr)

	//tryJwsOpt := echoapp_middlewares.JwsMiddlewaresOptions{
	//	Skipper:    middleware.DefaultSkipper,
	//	Jws:        app.MustGetJwsHelper(),
	//	IgnoreAuth: true,
	//}
	//tryJwsMiddleware := echoapp_middlewares.NewJwsMiddlewares(tryJwsOpt)

	normal := e.Group("/v1/order")
	normal.Use(corsMiddleware, limitMiddleware, companyMiddleware)
	orderCtl := controllers.NewOrderController(orderSvr)
	normal.GET("/getTicketByCode", orderCtl.GetTicketByCode)

	jwsAuth := e.Group("/v1/order")
	jwsOpt := echoapp_middlewares.JwsMiddlewaresOptions{
		Skipper: middleware.DefaultSkipper,
		Jws:     app.MustGetJwsHelper(),
	}
	jwsMiddleware := echoapp_middlewares.NewJwsMiddlewares(jwsOpt)
	jwsAuth.Use(corsMiddleware, jwsMiddleware, limitMiddleware, companyMiddleware)
	jwsAuth.GET("/getOrderList", orderCtl.GetOrderList)
	go func() {
		if err := e.Start(echoapp.ConfigOpts.OrderServer.Addr); err != nil {
			echoapp_util.DefaultLogger().WithError(err).Error("服务启动异常")
			os.Exit(-1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		echoapp_util.DefaultLogger().WithError(err).Error("服务关闭异常")
	}
}

// serverCmd represents the server command
var orderServerCmd = &cobra.Command{
	Use:   "order",
	Short: "商品订单服务",
	Long:  `商品订单服务`,
	Run: func(cmd *cobra.Command, args []string) {
		startOrderServer()
	},
}

func init() {
	rootCmd.AddCommand(orderServerCmd)
}
