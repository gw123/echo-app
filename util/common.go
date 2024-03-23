package echoapp_util

import (
	"strings"

	echoapp "github.com/gw123/echo-app"
)

//parse clientType by userAgent
func GetClientTypeByUA(ua string) string {
	ua = strings.ToLower(ua)
	if strings.Contains(ua, "micromessenger") {
		return echoapp.ClientWxOfficial
	} else {
		return echoapp.ClientWxMiniApp
	}
	return echoapp.ClientWap
}
