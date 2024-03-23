package bigdata

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/viper"

	"github.com/gw123/glog"

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
	e := echo.New()
	mode := echoapp.ConfigOpts.ApiVersion
	echoapp_util.DefaultLogger().Infof("开启站点服务 version %s", mode)

	e.HTTPErrorHandler = func(err error, ctx echo.Context) {
		ctx.JSON(http.StatusInternalServerError, map[string]string{"msg": err.Error()})
	}

	e.GET("/debug", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "1.9.24.2222")
	})
	//前端入口
	e.Static("/dev/public", echoapp.ConfigOpts.Asset.PublicRoot)
	e.Static("/", echoapp.ConfigOpts.Asset.PublicRoot)
	assetConfig := echoapp.ConfigOpts.Asset
	e.Renderer = echoapp_util.NewTemplateRenderer(assetConfig.ViewRoot, assetConfig.PublicHost, assetConfig.Version)

	origins := echoapp.ConfigOpts.SiteServer.Origins
	glog.DefaultLogger().Infof("origins [%+v]", origins)
	if len(origins) > 0 {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: origins,
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType,
				echo.HeaderAccept, "x-requested-with", "authorization", "ClientID", "x-csrf-token", "Access-Control-Allow-Credentials"},
		}))
	}

	glog.DefaultLogger().Infof("ConfigTraceAgentHostPort %s", viper.GetString(echoapp.ConfigTraceAgentHostPort))
	loggerMiddleware := echoapp_middlewares.NewLoggingMiddleware(echoapp_middlewares.LoggingMiddlewareConfig{
		Skipper: func(ctx echo.Context) bool {
			req := ctx.Request()
			return (req.RequestURI == "/" && req.Method == "HEAD") || (req.RequestURI == "/favicon.ico" && req.Method == "GET")
		},
		EnableTrace:        true,
		TraceAgentHostPort: viper.GetString(echoapp.ConfigTraceAgentHostPort),
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
	siteSvr := app.MustGetSiteService()
	userSvr := app.MustGetUserService()
	videoSvr := app.MustGetVideoService()

	companyCtl := controllers.NewCompanyController(comSvr)

	wechatSvr := app.MustGetWechatService()
	weChatMiddle := echoapp_middlewares.NewWechatAuthMiddlewares(
		middleware.DefaultSkipper,
		wechatSvr,
		userSvr,
	)

	tryJwsOpt := echoapp_middlewares.JwsMiddlewaresOptions{
		Skipper: middleware.DefaultSkipper,
		Jws:     app.MustGetJwsHelper(),
		//IgnoreAuth: true,
		IsTry: true,
	}
	tryJwsMiddle := echoapp_middlewares.NewJwsMiddlewares(tryJwsOpt)
	siteCtl := controllers.NewSiteController(comSvr, actSvr, siteSvr, wechatSvr, videoSvr, echoapp.ConfigOpts.Asset)

	wechatGroup := e.Group("/index-dev")
	wechatGroup.Use(companyMiddleware, tryJwsMiddle, weChatMiddle)
	wechatGroup.GET("/:com_id", siteCtl.Index)
	wechatGroup.GET("/:com_id/wxAuthCallBack", siteCtl.WxAuthCallBack)

	wechatGroup2 := e.Group("/index")
	wechatGroup2.Use(companyMiddleware, tryJwsMiddle, weChatMiddle)
	wechatGroup2.GET("/:com_id", siteCtl.Index)
	wechatGroup2.GET("/:com_id/wxAuthCallBack", siteCtl.WxAuthCallBack)

	e.GET("/index-dev/:com_id/video/:id", siteCtl.GetVideoDetail)
	e.GET("/index/:com_id/video/:id", siteCtl.GetVideoDetail)

	normal := e.Group("/" + mode + "/site/:com_id")
	normal.Use(companyMiddleware, limitMiddleware)
	//首页显示
	normal.GET("/getWxConfig", siteCtl.GetWxConfig, tryJwsMiddle)
	normal.GET("/wxMessage", siteCtl.WxMessage)
	normal.GET("/getBannerList", siteCtl.GetBannerList)
	normal.GET("/getIndexPageBanners", siteCtl.GetIndexPageBanners)

	normal.GET("/getNotifyList", siteCtl.GetNotifyList)
	normal.GET("/getNotifyDetail", siteCtl.GetNotifyDetail)
	normal.GET("/getActivityList", siteCtl.GetActivityList)
	normal.GET("/getActivityDetail", siteCtl.GetActivityDetail)
	normal.GET("/getNavList", siteCtl.GetQuickNav)
	normal.GET("/getCompany", companyCtl.GetCompanyInfo)
	normal.GET("/getVideoList", siteCtl.GetVideoList)
	//
	normal.POST("/sendMailCode", siteCtl.SendMailCode)
	normal.POST("/sendRawMail", siteCtl.SendRawMail)

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
var serverCmd = &cobra.Command{
	Use:   "bigdata",
	Short: "大数据",
	Long:  `大数据`,
	Run: func(cmd *cobra.Command, args []string) {
		startSiteServer()
	},
}
