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
	"github.com/gw123/echo-app/app"
	"github.com/gw123/echo-app/controllers"
	echoapp_middlewares "github.com/gw123/echo-app/middlewares"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/cobra"
)

func startFileServer() {
	echoapp_util.DefaultLogger().Info("开启HTTP服务")
	e := echo.New()
	e.HTTPErrorHandler = func(err error, ctx echo.Context) {
		ctx.JSON(http.StatusInternalServerError, map[string]string{"msg": err.Error()})
	}

	if echoapp.ConfigOpts.Asset.PublicRoot != "" {
		e.Static("/", echoapp.ConfigOpts.Asset.PublicRoot)
	}

	assetConfig := echoapp.ConfigOpts.Asset
	e.Renderer = echoapp_util.NewTemplateRenderer(assetConfig.ViewRoot, assetConfig.PublicHost, assetConfig.Version)

	origins := echoapp.ConfigOpts.FileServer.Origins
	if len(origins) > 0 {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: origins,
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAcceptEncoding,
				"Accept-Language", "Referer", "Connection", "ClientID",
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
	//e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
	//	StackSize: 1 << 10, // 1 KB
	//}))

	//Actions
	goodSvr := app.MustGetGoodsService()
	resourceSvc := app.MustGetResourceService()
	resourceCtl := controllers.NewResourceController(resourceSvc, goodSvr)
	e.GET("/v1/file/ping", func(c echo.Context) error {
		if app.App.IsHealth {
			c.HTML(http.StatusOK, "pong")
		} else {
			c.HTML(http.StatusInternalServerError, "not health")
		}
		return nil
	})

	tryJwsAuthGroup := e.Group("/v1/file")

	tryJwsOpt := echoapp_middlewares.JwsMiddlewaresOptions{
		Skipper:    middleware.DefaultSkipper,
		Jws:        app.MustGetJwsHelper(),
		IgnoreAuth: true,
	}
	tryJwsMiddleware := echoapp_middlewares.NewJwsMiddlewares(tryJwsOpt)
	tryJwsAuthGroup.Use(tryJwsMiddleware)
	tryJwsAuthGroup.POST("/uploadImage", resourceCtl.UploadImage)

	jwsAuth := e.Group("/v1/file")
	jwsOpt := echoapp_middlewares.JwsMiddlewaresOptions{
		Skipper: middleware.DefaultSkipper,
		Jws:     app.MustGetJwsHelper(),
	}
	jwsMiddleware := echoapp_middlewares.NewJwsMiddlewares(jwsOpt)
	jwsAuth.Use(jwsMiddleware)
	//jwsAuth.POST("/saveReource", resourceCtl.SaveResource)
	jwsAuth.GET("/getResourceById", resourceCtl.GetResourceById)
	jwsAuth.GET("/getResourcesByTagId", resourceCtl.GetResourcesByTagId)
	jwsAuth.GET("/getUserPaymentResources", resourceCtl.GetUserPaymentResources)
	jwsAuth.POST("/uploadResource", resourceCtl.UploadResource)

	jwsAuth.GET("/getResourceList", resourceCtl.GetResourceList)
	jwsAuth.GET("/getResourceByName", resourceCtl.GetResourceByName)
	jwsAuth.GET("/downloadFile", resourceCtl.DownloadResource)
	jwsAuth.GET("/getSelfResources", resourceCtl.GetSelfResources)
	jwsAuth.GET("/getUserPaymentResources", resourceCtl.GetUserPaymentResources)
	go func() {
		if err := e.Start(echoapp.ConfigOpts.FileServer.Addr); err != nil {
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
var fileServerCmd = &cobra.Command{
	Use:   "file",
	Short: "文件服务",
	Long:  `文件服务`,
	Run: func(cmd *cobra.Command, args []string) {
		startFileServer()
	},
}

func init() {
	rootCmd.AddCommand(fileServerCmd)
}
