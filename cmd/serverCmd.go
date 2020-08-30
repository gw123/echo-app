// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/controllers"
	echoapp_middlewares "github.com/gw123/echo-app/middlewares"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/cobra"
)

func startHttp() {
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

	origins := echoapp.ConfigOpts.Server.Origins
	if len(origins) > 0 {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: origins,
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType,
				echo.HeaderAccept, "x-requested-with", "authorization", "x-csrf-token"},
		}))
	}

	loggerMiddleware := echoapp_middlewares.NewLoggingMiddleware(echoapp_middlewares.LoggingMiddlewareConfig{
		Skipper: func(ctx echo.Context) bool {
			req := ctx.Request()
			return (req.RequestURI == "/" && req.Method == "HEAD") || (req.RequestURI == "/favicon.ico" && req.Method == "GET")
		},
	})
	e.Use(loggerMiddleware)
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
	}))

	//Actions
	//	usrSvr := app.MustUserService()
	//	resourceSvr := app.MustGetResourceService()
	//orderSvr := app.MustGetOrderService()
	//exampleController := controllers.ExampleController{}
	// qrcodeController := controllers.NewQrcodeController()
	areaCtr := controllers.NewAreaController()
	smsCtr := controllers.NewSmsController()
	//	userCtl := controllers.NewUserController(usrSvr)
	//	resourceCtl := controllers.NewResourceController(resourceSvr)
	//orderCtl := controllers.NewOrderController(orderSvr)
	// e.GET("/index", exampleController.Index)
	// e.GET("/getQrcode", qrcodeController.GetQrcode)
	e.GET("/getAreaMap", areaCtr.GetAreaMap)
	e.GET("/getAreaArray", areaCtr.GetAreaArray)
	e.POST("/sendMessage", smsCtr.SendMessageByToken)

	// e.POST("/addUserScore", userCtl.AddUserScore)
	// e.POST("/subUserScore", userCtl.SubUserScore)
	// e.POST("/addUser", userCtl.Register)
	// e.POST("/login", userCtl.Login)
	// e.POST("/addUserRole", userCtl.AddUserRoles)
	// e.POST("/addPermission", userCtl.AddPermissions)
	// e.POST("/roleHasPermission", userCtl.RoleHasPermission)

	// e.POST("/saveReource", resourceCtl.SaveResource)
	// e.POST("/getResourcebyId", resourceCtl.GetResourceById)
	// e.POST("/getResourcesbyTagId", resourceCtl.GetResourcesByTagId)
	// e.POST("/getUserPamentResources", resourceCtl.GetUserPaymentResources)
	// e.POST("/uploadResource", resourceCtl.UploadResource)
	// e.GET("/getResourceList", resourceCtl.GetResourceList)
	// //e.GET("/getResourcebyPath", resourceCtl.GetResourceByName)
	// e.GET("/getResourcebyName", resourceCtl.GetResourceByName)
	// e.GET("/downloadFile", resourceCtl.DownloadResource)

	// e.POST("/placeOrder", orderCtl.UniPreOrder)
	// e.POST("/getOrderbyId", orderCtl.GetOrdereById)
	// e.POST("/getUserPaymentOrder", orderCtl.GetUserPaymentOrder)
	// e.POST("/getOrderList", orderCtl.GetOrderList)
	go func() {
		if err := e.Start(echoapp.ConfigOpts.Server.Addr); err != nil {
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
	Use:   "server",
	Short: "服务",
	Long:  `测试服务`,
	Run: func(cmd *cobra.Command, args []string) {
		startHttp()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
