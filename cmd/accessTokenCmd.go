package cmd

import (
	"encoding/json"
	"fmt"

	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type AccessTokenOption struct {
	GrantType string `json:"grant_type"`
	AppId     string `json:"appid"`
	Secret    string `json:"secret"`
}
type ResponseToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func getAccessToken() error {
	echoapp_util.DefaultJsonLogger().Infof("获取微信AccessToken")
	getTokenRequest := AccessTokenOption{
		GrantType: "client_credential",
		AppId:     "wx97a1f596e4b87a83",
		Secret:    "ea389cb3939a61965a6cd5aba79e09e0",
	}
	baseurl := "https://api.weixin.qq.com/cgi-bin/token"
	url := baseurl + fmt.Sprintf("?grant_type=%s&appid=%s&secret=%s",
		getTokenRequest.GrantType,
		getTokenRequest.AppId,
		getTokenRequest.Secret)
	res, err := echoapp_util.DoHttpRequest(url, "GET", nil)
	if err != nil {
		return errors.Wrap(err, "echoapp_util.DoHttpRequest")
	}
	responseToken := &ResponseToken{}
	if err = json.Unmarshal(res, responseToken); err != nil {
		return errors.Wrap(err, "doReportHttpRequest")
	}
	echoapp_util.DefaultLogger().Info(string(res))
	return nil
}

var accessTokenCmd = &cobra.Command{
	Use:   "wx-access-token",
	Short: "服务",
	Long:  `测试服务`,
	Run: func(cmd *cobra.Command, args []string) {
		getAccessToken()
	},
}

func init() {
	rootCmd.AddCommand(accessTokenCmd)
}
