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
	"fmt"
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/controllers"
	echoapp_middlewares "github.com/gw123/echo-app/middlewares"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/cobra"
	"net/http"
)

func startHttp() {
	echoapp_util.DefaultLogger().Infof("%+v", echoapp.Config)
	e := echo.New()
	e.HTTPErrorHandler = func(err error, ctx echo.Context) {
		ctx.JSON(http.StatusInternalServerError, map[string]string{"msg": err.Error()})
	}

	if echoapp.Config.Asset.PublicRoot != "" {
		e.Static("/", echoapp.Config.Asset.PublicRoot)
	}

	assetConfig := echoapp.Config.Asset
	fmt.Println(assetConfig)
	e.Renderer = echoapp_util.NewTemplateRenderer(assetConfig.ViewRoot, assetConfig.PublicHost, assetConfig.Version)

	origins := echoapp.Config.Origins
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
	exampleController := controllers.ExampleController{}
	qrcodeController := controllers.NewQrcodeController()
	areaCtl := controllers.NewAreaController()
	e.GET("/index", exampleController.Index)
	e.GET("/getQrcode", qrcodeController.GetQrcode)
	e.GET("/getAreaList", areaCtl.GetAreaList)
	//
	if err := e.Start(echoapp.Config.Env.Addr); err != nil {
		echoapp_util.DefaultLogger().WithError(err).Error("服务启动异常")
	}
	ctx := context.Background()
	if err := e.Shutdown(ctx); err != nil {
		echoapp_util.DefaultLogger().WithError(err).Error("服务关闭异常")
	}
}

// serverCmd represents the server command
var innerServerCmd = &cobra.Command{
	Use:   "server",
	Short: "服务",
	Long:  `测试服务`,
	Run: func(cmd *cobra.Command, args []string) {
		startHttp()
	},
}

func init() {
	rootCmd.AddCommand(innerServerCmd)
}
