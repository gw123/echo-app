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

func startUserServer() {
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

	origins := echoapp.ConfigOpts.UserServer.Origins
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

	companySvr := app.MustGetCompanyService()
	limitMiddleware := echoapp_middlewares.NewLimitMiddlewares(middleware.DefaultSkipper, 100, 200)
	companyMiddleware := echoapp_middlewares.NewCompanyMiddlewares(middleware.DefaultSkipper, companySvr)
	usrSvr := app.MustGetUserService()
	goodsSvr := app.MustGetGoodsService()
	smsSvr := app.MustGetSmsService()
	wechatSvr := app.MustGetWechatService()
	userCtl := controllers.NewUserController(usrSvr, goodsSvr, smsSvr, wechatSvr)
	mode := echoapp.ConfigOpts.ApiVersion
	normal := e.Group("/" + mode + "/user/:com_id")
	tryJwsOpt := echoapp_middlewares.JwsMiddlewaresOptions{
		Skipper:    middleware.DefaultSkipper,
		Jws:        app.MustGetJwsHelper(),
		IgnoreAuth: true,
		//MockUserId: 58,
	}
	normal.Use(companyMiddleware, echoapp_middlewares.NewJwsMiddlewares(tryJwsOpt))
	//登录
	normal.POST("/login", userCtl.Login)
	normal.POST("/register", userCtl.Register)
	normal.POST("/logout", userCtl.Logout)
	normal.POST("/sendVerifyCodeSms", userCtl.SendVerifyCodeSms)
	//normal.POST("/checkVerifyCode", userCtl.CheckVerifyCode)
	//

	normal.GET("/getVerifyPic", userCtl.GetVerifyPic)

	//jwsAuth := e.Group("/v1/user")
	jwsAuth := e.Group("/" + mode + "/user/:com_id")
	jwsOpt := echoapp_middlewares.JwsMiddlewaresOptions{
		Skipper:    middleware.DefaultSkipper,
		Jws:        app.MustGetJwsHelper(),
		//IgnoreAuth: true,
		//MockUserId: 58,
	}
	jwsMiddleware := echoapp_middlewares.NewJwsMiddlewares(jwsOpt)
	userMiddleware := echoapp_middlewares.NewUserMiddlewares(middleware.DefaultSkipper, usrSvr)
	jwsAuth.Use(jwsMiddleware, limitMiddleware, userMiddleware)
	jwsAuth.POST("/changeUserScore", userCtl.AddUserScore)
	jwsAuth.POST("/jscode2session", userCtl.Jscode2session)
	jwsAuth.GET("/getUserInfo", userCtl.GetUserInfo)
	//roles
	jwsAuth.POST("/getUserRoles", userCtl.GetUserRoles)
	jwsAuth.POST("/checkHasRoles", userCtl.CheckHasRoles)
	//addressList
	jwsAuth.GET("/getUserAddressList", userCtl.GetUserAddressList)
	jwsAuth.POST("/createUserAddress", userCtl.CreateUserAddress)
	jwsAuth.POST("/updateUserAddress", userCtl.UpdateUserAddress)
	jwsAuth.POST("/delUserAddress", userCtl.DelUserAddress)
	jwsAuth.GET("/getUserDefaultAddress", userCtl.GetUserDefaultAddress)
	//collection
	jwsAuth.GET("/getUserCollectionList", userCtl.GetUserCollectionList)
	jwsAuth.POST("/addUserCollection", userCtl.AddUserCollection)
	jwsAuth.POST("/delUserCollection", userCtl.DelUserCollection)
	jwsAuth.POST("/isCollect", userCtl.IsCollect)
	//history
	jwsAuth.POST("/addUserHistory", userCtl.AddUserHistory)
	jwsAuth.GET("/getUserHistory", userCtl.GetUserHistoryList)
	jwsAuth.GET("/getUserBrowseLeaderboard", userCtl.GetUserBrowseLeaderboard)
	go func() {
		if err := e.Start(echoapp.ConfigOpts.UserServer.Addr); err != nil {
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
var userServerCmd = &cobra.Command{
	Use:   "user",
	Short: "用户服务",
	Long:  `用户服务`,
	Run: func(cmd *cobra.Command, args []string) {
		startUserServer()
	},
}

func init() {
	rootCmd.AddCommand(userServerCmd)
}
