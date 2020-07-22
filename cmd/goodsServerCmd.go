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

func startGoodsServer() {
	echoapp_util.DefaultLogger().Info("开启GoodsHTTP服务")
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

	corsMiddleware := middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: echoapp.ConfigOpts.GoodsServer.Origins,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, "ClientID",
			echo.HeaderAccept, "x-requested-with", "authorization", "x-csrf-token", "Access-Control-Allow-Credentials"},
	})

	loggerMiddleware := echoapp_middlewares.NewLoggingMiddleware(echoapp_middlewares.LoggingMiddlewareConfig{
		Skipper: func(ctx echo.Context) bool {
			req := ctx.Request()
			return (req.RequestURI == "/" && req.Method == "HEAD") || (req.RequestURI == "/favicon.ico" && req.Method == "GET")
		},
	})
	e.Use(corsMiddleware, loggerMiddleware)
	//e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
	//	StackSize: 1 << 10, // 1 KB
	//}))

	//Actions
	companySvr := app.MustGetCompanyService()
	goodsSvr := app.MustGetGoodsService()
	limitMiddleware := echoapp_middlewares.NewLimitMiddlewares(middleware.DefaultSkipper, 100, 200)
	companyMiddleware := echoapp_middlewares.NewCompanyMiddlewares(middleware.DefaultSkipper, companySvr)

	tryJwsOpt := echoapp_middlewares.JwsMiddlewaresOptions{
		Skipper:    middleware.DefaultSkipper,
		Jws:        app.MustGetJwsHelper(),
		IgnoreAuth: true,
	}
	tryJwsMiddleware := echoapp_middlewares.NewJwsMiddlewares(tryJwsOpt)
	mode := echoapp.ConfigOpts.ApiVersion
	normal := e.Group("/" + mode + "/goods/:com_id")
	normal.Use( limitMiddleware, companyMiddleware, tryJwsMiddleware)

	goodsCtl := controllers.NewGoodsController(goodsSvr)

	normal.GET("/getGoodsList", goodsCtl.GetGoodsList)
	normal.GET("/getRecommendGoodsList", goodsCtl.GetRecommendGoodsList)
	normal.GET("/getGoodsListByTagId", goodsCtl.GetTagGoodsList)
	normal.GET("/getGoodsTags", goodsCtl.GetGoodsTags)
	normal.GET("/getGoodsDetail", goodsCtl.GetGoodsInfo)

	//cart
	jwsAuth := e.Group("/" + mode + "/goods/:com_id")
	jwsMiddleware := echoapp_middlewares.NewJwsMiddlewares(echoapp_middlewares.JwsMiddlewaresOptions{
		Skipper: middleware.DefaultSkipper,
		Jws:     app.MustGetJwsHelper(),
	})
	jwsAuth.Use(jwsMiddleware, limitMiddleware, companyMiddleware)
	jwsAuth.GET("/getCartGoodsList", goodsCtl.GetCartGoodsList)
	jwsAuth.POST("/addCartGoods", goodsCtl.AddCartGoods)
	jwsAuth.POST("/delCartGoods", goodsCtl.DelCartGoods)
	jwsAuth.POST("/clearCart", goodsCtl.ClearCart)
	jwsAuth.POST("/updateCartGoods", goodsCtl.UpdateCartGoods)

	go func() {
		if err := e.Start(echoapp.ConfigOpts.GoodsServer.Addr); err != nil {
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
var goodsServerCmd = &cobra.Command{
	Use:   "goods",
	Short: "商品服务",
	Long:  `商品服务`,
	Run: func(cmd *cobra.Command, args []string) {
		startGoodsServer()
	},
}

func init() {
	rootCmd.AddCommand(goodsServerCmd)
}
