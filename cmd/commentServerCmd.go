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

func startCommentServer() {
	echoapp_util.DefaultLogger().Info("开启HTTP服务")
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

	origins := echoapp.ConfigOpts.CommentServer.Origins
	if len(origins) > 0 {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: origins,
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType,
				echo.HeaderAccept, "x-requested-with", "authorization", "x-csrf-token", "ClientID", "Access-Control-Allow-Credentials"},
		}))
	}

	loggerMiddleware := echoapp_middlewares.NewLoggingMiddleware(echoapp_middlewares.LoggingMiddlewareConfig{
		Skipper: func(ctx echo.Context) bool {
			req := ctx.Request()
			return (req.RequestURI == "/" && req.Method == "HEAD") || (req.RequestURI == "/favicon.ico" && req.Method == "GET")
		},
		Logger: glog.JsonEntry(),
	})
	e.Use(loggerMiddleware)
	//e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
	//	StackSize: 1 << 10, // 1 KB
	//}))

	//Actions
	commentSvr := app.MustGetCommentService()
	//goodsSvr := app.MustGetGoodsService()
	//resourceSvr := app.MustGetResourceService()
	limitMiddleware := echoapp_middlewares.NewLimitMiddlewares(middleware.DefaultSkipper, 100, 200)
	//companyMiddleware := echoapp_middlewares.NewCompanyMiddlewares(middleware.DefaultSkipper,goodsSvr)

	tryJwsOpt := echoapp_middlewares.JwsMiddlewaresOptions{
		Skipper:    middleware.DefaultSkipper,
		Jws:        app.MustGetJwsHelper(),
		IgnoreAuth: true,
	}
	tryJwsMiddleware := echoapp_middlewares.NewJwsMiddlewares(tryJwsOpt)
	//resourceCtl := controllers.NewResourceController(resourceSvr, goodsSvr)
	//
	//callback := e.Group("/v1/goods-api")
	//callback.POST("/uploadCallback", resourceCtl.UploadCallback)
	//
	mode := echoapp.ConfigOpts.ApiVersion
	normal := e.Group("/" + mode + "/comment/:com_id")
	//normal := e.Group("/v1/comment")
	normal.Use(limitMiddleware, tryJwsMiddleware)

	commentCtl := controllers.NewCommentController(commentSvr)
	normal.POST("/submitComment", commentCtl.SaveComment)
	normal.GET("/getCommentList", commentCtl.GetCommentList)
	normal.GET("/getGoodsCommentNum", commentCtl.GetGoodsCommentNum)
	normal.GET("/getSubCommentList", commentCtl.GetSubCommentList)
	normal.GET("/upComment", commentCtl.ThumbUpComment)
	go func() {
		if err := e.Start(echoapp.ConfigOpts.CommentServer.Addr); err != nil {
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
var commentServerCmd = &cobra.Command{
	Use:   "comment",
	Short: "评论服务",
	Long:  `评论服务`,
	Run: func(cmd *cobra.Command, args []string) {
		startCommentServer()
	},
}

func init() {
	rootCmd.AddCommand(commentServerCmd)
}
