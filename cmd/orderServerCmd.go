package cmd

import (
	"context"
	"github.com/gw123/glog"
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
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, "ClientID",
			echo.HeaderAccept, "x-requested-with", "authorization", "x-csrf-token", "Access-Control-Allow-Credentials"},
	})

	loggerMiddleware := echoapp_middlewares.NewLoggingMiddleware(echoapp_middlewares.LoggingMiddlewareConfig{
		Skipper: func(ctx echo.Context) bool {
			req := ctx.Request()
			return (req.RequestURI == "/" && req.Method == "HEAD") || (req.RequestURI == "/favicon.ico" && req.Method == "GET")
		},
		Logger: glog.JsonEntry(),
	})
	e.Use(corsMiddleware, loggerMiddleware)
	//e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
	//	StackSize: 1 << 10, // 1 KB
	//}))

	//Actions
	orderSvr := app.MustGetOrderService()
	companySvr := app.MustGetCompanyService()
	limitMiddleware := echoapp_middlewares.NewLimitMiddlewares(middleware.DefaultSkipper, 100, 200)
	companyMiddleware := echoapp_middlewares.NewCompanyMiddlewares(middleware.DefaultSkipper, companySvr)

	tryJwsMiddleware := echoapp_middlewares.NewJwsMiddlewares(echoapp_middlewares.JwsMiddlewaresOptions{
		Skipper:    middleware.DefaultSkipper,
		Jws:        app.MustGetJwsHelper(),
		IgnoreAuth: true,
	})
	mode := echoapp.ConfigOpts.ApiVersion
	normal := e.Group("/" + mode + "/order/:com_id")
	normal.Use(limitMiddleware, companyMiddleware, tryJwsMiddleware)
	userSvr := app.MustGetUserService()
	orderCtl := controllers.NewOrderController(orderSvr, userSvr)
	normal.GET("/getTicketByCode", orderCtl.GetTicketByCode)
	//微信支付回调
	normal.POST("/wxPayCallback", orderCtl.WxPayCallback)
	normal.POST("/wxRefundCallback", orderCtl.WxRefundCallback)
	jwsAuth := e.Group("/" + mode + "/order/:com_id")
	jwsMiddleware := echoapp_middlewares.NewJwsMiddlewares(echoapp_middlewares.JwsMiddlewaresOptions{
		Skipper: middleware.DefaultSkipper,
		Jws:     app.MustGetJwsHelper(),
	})
	userMiddle := echoapp_middlewares.NewUserMiddlewares(middleware.DefaultSkipper, userSvr)
	jwsAuth.Use(jwsMiddleware, userMiddle, limitMiddleware, companyMiddleware)
	jwsAuth.GET("/getOrderList", orderCtl.GetOrderList)
	jwsAuth.GET("/getOrderDetail", orderCtl.GetOrderDetail)
	jwsAuth.GET("/getOrderStatistics", orderCtl.GetOrderStatistics)
	jwsAuth.POST("/preOrder", orderCtl.PreOrder)
	jwsAuth.POST("/queryOrder", orderCtl.QueryOrder)
	jwsAuth.POST("/cancelOrder", orderCtl.CancelOrder)
	jwsAuth.POST("/refund", orderCtl.Refund)
	jwsAuth.POST("/queryRefund", orderCtl.QueryRefund)

	//ticket
	jwsAuth.GET("/checkTicketByStaff", orderCtl.CheckTicketByStaff)
	jwsAuth.GET("/checkTicketBySelf", orderCtl.CheckTicketBySelf)
	jwsAuth.GET("/checkTicketList", orderCtl.CheckTicketList)
	jwsAuth.GET("/getTicketList", orderCtl.GetTicketList)
	jwsAuth.GET("/getTicketDetail", orderCtl.GetTicketDetail)
	jwsAuth.GET("/fetchThirdTicket", orderCtl.FetchThirdTicket)

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
	Short: "商品服务",
	Long:  `商品服务`,
	Run: func(cmd *cobra.Command, args []string) {
		startOrderServer()
	},
}

func init() {
	rootCmd.AddCommand(orderServerCmd)
}
