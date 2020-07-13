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

func startSiteServer() {
	echoapp_util.DefaultLogger().Info("开启站点服务")
	e := echo.New()
	e.HTTPErrorHandler = func(err error, ctx echo.Context) {
		ctx.JSON(http.StatusInternalServerError, map[string]string{"msg": err.Error()})
	}
	//前端入口
	e.Static("/", echoapp.ConfigOpts.Asset.PublicRoot+"/m")
	assetConfig := echoapp.ConfigOpts.Asset
	e.Renderer = echoapp_util.NewTemplateRenderer(assetConfig.ViewRoot, assetConfig.PublicHost, assetConfig.Version)

	origins := echoapp.ConfigOpts.SiteServer.Origins
	if len(origins) > 0 {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: origins,
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType,
				echo.HeaderAccept, "x-requested-with", "authorization", "ClientID", "x-csrf-token", "Access-Control-Allow-Credentials"},
		}))
	}

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
	companySvr := app.MustGetCompanyService()
	limitMiddleware := echoapp_middlewares.NewLimitMiddlewares(middleware.DefaultSkipper, 100, 200)
	companyMiddleware := echoapp_middlewares.NewCompanyMiddlewares(middleware.DefaultSkipper, companySvr)
	comSvr := app.MustGetCompanyService()
	actSvr := app.MustGetActivityService()
	siteCtl := controllers.NewSiteController(comSvr, actSvr)
	companyCtl := controllers.NewCompanyController(comSvr)

	mode := echoapp.ConfigOpts.ApiVersion
	normal := e.Group("/" + mode + "/site/:com_id")
	tryJwsOpt := echoapp_middlewares.JwsMiddlewaresOptions{
		Skipper:    middleware.DefaultSkipper,
		Jws:        app.MustGetJwsHelper(),
		IgnoreAuth: true,
	}
	normal.Use(companyMiddleware, limitMiddleware, echoapp_middlewares.NewJwsMiddlewares(tryJwsOpt))
	//首页显示
	normal.GET("/getBannerList", siteCtl.GetBannerList)
	normal.GET("/getNotifyList", siteCtl.GetNotifyList)
	normal.GET("/getNotifyDetail", siteCtl.GetNotifyDetail)
	normal.GET("/getActivityList", siteCtl.GetActivityList)
	normal.GET("/getActivityDetail", siteCtl.GetActivityDetail)
	normal.GET("/getNavList", siteCtl.GetQuickNav)
	normal.GET("/getCompany", companyCtl.GetCompanyInfo)

	go func() {
		if err := e.Start(echoapp.ConfigOpts.SiteServer.Addr); err != nil {
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
var siteServerCmd = &cobra.Command{
	Use:   "site",
	Short: "站点服务",
	Long:  `站点服务`,
	Run: func(cmd *cobra.Command, args []string) {
		startSiteServer()
	},
}

func init() {
	rootCmd.AddCommand(siteServerCmd)
}
