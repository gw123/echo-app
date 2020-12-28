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
	"crypto/sha1"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"time"

	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"github.com/spf13/cobra"
)

func startWXHttp() {
	echoapp_util.DefaultLogger().Info("开启公众号服务")
	e := echo.New()
	e.HTTPErrorHandler = func(err error, ctx echo.Context) {
		ctx.JSON(http.StatusInternalServerError, map[string]string{"msg": err.Error()})
	}
	e.GET("/", func(c echo.Context) error {
		signature := c.QueryParam("signature")
		timestamp := c.QueryParam("timestamp")
		nonce := c.QueryParam("nonce")
		echostr := c.QueryParam("echostr")
		token := "Gaohng"
		tmpArr := []string{token, timestamp, nonce}
		sort.Strings(tmpArr)
		//sort.Sort(sort.StringSlice(tmpArr))
		var build strings.Builder
		for _, val := range tmpArr {
			build.WriteString(val)
		}
		tmpStr := build.String()
		//sha1加密
		h := sha1.New()
		h.Write([]byte(tmpStr))
		bs := h.Sum(nil)
		if string(bs) == signature {
			return c.String(http.StatusOK, echostr)
		}
		return nil
	})

	go func() {
		if err := e.Start("0.0.0.0:80"); err != nil {
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
var wxServerCmd = &cobra.Command{
	Use:   "wxOfficialAccount",
	Short: "服务",
	Long:  `测试服务`,
	Run: func(cmd *cobra.Command, args []string) {
		startWXHttp()
	},
}

func init() {
	rootCmd.AddCommand(wxServerCmd)
}
