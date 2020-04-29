package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/app"

	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/spf13/cobra"
)

//监控目录
func watchDir(dir string) {
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	//通过Walk来遍历目录下的所有子目录
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		//这里判断是否为目录，只需监控目录即可
		//目录下的文件也在监控范围内，不需要我们一个一个加
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}

			err = watch.Add(path)
			if err != nil {
				return err
			}
			fmt.Println("监控 : ", path)
		}

		return nil
	})

	resourceSvc := app.MustGetResService()
	go func() {
		for {
			select {
			case ev := <-watch.Events:
				{
					if ev.Op&fsnotify.Create == fsnotify.Create {
						fmt.Println("创建文件 : ", ev.Name)

						//这里获取新创建文件的信息，如果是目录，则加入监控中
						fi, err := os.Stat(ev.Name)
						if err == nil && fi.IsDir() {
							watch.Add(ev.Name)
							fmt.Println("添加监控 : ", ev.Name)
						} else {
							Md5file, err := resourceSvc.Md5SumFile(ev.Name)
							if err != nil {
								fmt.Println(err)
							}
							resourceSvc.SaveResource(&echoapp.Resource{
								Path:    ev.Name,
								Md5File: Md5file,
								Type:    path.Ext(ev.Name),
							})
						}
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						fmt.Println("写入文件 : ", ev.Name)
						Md5file, err := resourceSvc.Md5SumFile(ev.Name)
						if err != nil {
							fmt.Println(err)
						}
						res, err := resourceSvc.GetResourceByPath(ev.Name)
						if err != nil {
							fmt.Println(err)
						}
						res.Md5File = Md5file
						resourceSvc.ModifyResource(res)
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {

						//如果删除文件是目录，则移除监控
						fi, err := os.Stat(ev.Name)
						if err == nil && fi.IsDir() {
							watch.Remove(ev.Name)
							fmt.Println("删除监控 : ", ev.Name)
						} else {
							res, err := resourceSvc.GetResourceByPath(ev.Name)
							if err != nil {
								panic(err)
							}
							resourceSvc.DeleteResource(res)
						}
					}
					if ev.Op&fsnotify.Rename == fsnotify.Rename {
						fmt.Println("重命名文件(删除文件) : ", ev.Name)
						res, err := resourceSvc.GetResourceByPath(ev.Name)

						if err != nil {
							fmt.Println(err)
						}
						if err := resourceSvc.DeleteResource(res); err != nil {
							fmt.Println(err)
						}

						//如果重命名文件是目录，则移除监控
						//注意这里无法使用os.Stat来判断是否是目录了
						//因为重命名后，go已经无法找到原文件来获取信息了
						//所以这里就简单粗爆的直接remove好了
						watch.Remove(ev.Name)
					}

					if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
						fmt.Println("修改权限 : ", ev.Name)
					}
				}

			case err := <-watch.Errors:
				{
					fmt.Println("error : ", err)
					return
				}
			}
		}
	}()
}

func startWatcher() {

	echoapp_util.DefaultLogger().Infof("开始监控服务")
	watchDir(echoapp.ConfigOpts.Asset.WatchRoot)
	//select {}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

}

var watcherCmd = &cobra.Command{
	Use:   "watch",
	Short: "监听",
	Long:  "监听目录",
	Run: func(cmd *cobra.Command, args []string) {
		startWatcher()
	},
}

func init() {
	rootCmd.AddCommand(watcherCmd)
}
