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
	"fmt"
	"log"
	"os"
	"strings"

	echoapp "github.com/gw123/echo-app"

	"github.com/gw123/glog"

	"github.com/spf13/cobra"
)

var (
	etcdEndpoints []string
	cfgFile       string
	cfgType       string
	etcdUsername  string
	etcdPassword  string
	etcdNamespace string
	etcdPath      string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "echoapp",
	Short: "echoapp for echo ,esay to develop go web application",
	Long:  `echoapp for echo ,esay to develop go web application`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	glog.SetDefaultLoggerFormatter(glog.GetDefaultJsonLoggerFormatter())
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	RootCmd.PersistentFlags().StringVar(&cfgType, "config-type", "file", "local|etcd")

	RootCmd.PersistentFlags().StringSliceVar(&etcdEndpoints, "etcd-endpoints", []string{}, "endpoints")
	RootCmd.PersistentFlags().StringVar(&etcdUsername, "etcd-username", "", "username")
	RootCmd.PersistentFlags().StringVar(&etcdPassword, "etcd-password", "", "password")
	RootCmd.PersistentFlags().StringVar(&etcdNamespace, "etcd-namespace", "", "namespace")
	RootCmd.PersistentFlags().StringVar(&etcdPath, "etcd-path", "", "etcdPath config data path")

	if len(etcdEndpoints) == 0 {
		etcdEndpoints = strings.Split(os.Getenv("ETCD_ENDPOINTS"), ",")
		if len(etcdEndpoints) == 0 {
			etcdEndpoints = append(etcdEndpoints, "http://127.0.0.1:2379")
		}
	}

	if etcdUsername == "" {
		etcdUsername = os.Getenv("ETCD_USERNAME")
	}

	if etcdPassword == "" {
		etcdPassword = os.Getenv("ETCD_PASSWORD")
	}

	if etcdNamespace == "" {
		etcdNamespace = os.Getenv("ETCD_NAMESPACE")
		if etcdNamespace == "" {
			etcdNamespace = "/xyt"
		}
	}

	if etcdPath == "" {
		etcdPath = os.Getenv("ETCD_PATH")
		if etcdPath == "" {
			etcdPath = "/config.yaml"
		}
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgType == "file" {
		glog.Info("load config from file")
		echoapp.LoadFromFile(cfgFile)
	} else {
		glog.Infof("load config from etcd addr:%s ,username:%s", strings.Join(etcdEndpoints, ","), etcdUsername)
		glog.Infof("load config from etcd: namespace:%s ,path:%s", etcdNamespace, etcdPath)
		echoapp.LoadFromEtcd(etcdEndpoints, etcdNamespace, etcdPath, etcdUsername, etcdPassword)
	}
}

func handleInitError(module string, err error) {
	if err == nil {
		return
	}
	log.Fatalf("init %s failed, err: %s", module, err)
}
