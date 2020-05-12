package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"regexp"

	"github.com/fsnotify/fsnotify"
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/app"
	"github.com/pkg/errors"

	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/spf13/cobra"
)

//监控目录
func copy(src, dst string) (string, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return "", err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return "", fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer destination.Close()
	if _, err = io.Copy(destination, source); err != nil {
		return "", err
	}
	return "", nil
}

func doHttpRequest(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "DoRequest->http.NewRequest")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "DoRequest->Do")
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "DoRequest->ReadAll")
	}
	if res.StatusCode == http.StatusOK {
		return data, errors.Wrapf(err, "StatusCode:%d", res.StatusCode)
	}
	return data, nil
}

func getPPTCoverUrl(pptUrl string) ([]string, error) {
	testMap := echoapp.ConfigOpts.TestMap
	for key, options := range testMap {
		echoapp_util.DefaultLogger().Infof("访问%s,com_id:%d", key, options.ComId)
		//ppturl :=
		url := options.BaseUrl + "onlinePreview" + "?url=" + pptUrl
		data, err := doHttpRequest(url)
		if err != nil {
			return nil, errors.Wrap(err, "doHttpRequest")
		}
		strdata := string(data)

		reg := regexp.MustCompile(`\<img .*?title=\"查看大图\" .*?data-src=\"(\S*)\" .*?\>`)
		if reg == nil {
			return nil, errors.Wrap(err, "regexp.MustCompile err")
		}
		res := reg.FindAllStringSubmatch(strdata, -1)
		urls := make([]string, 0)
		for _, text := range res {
			url := text[1]
			urls = append(urls, url)
		}
		return urls, nil
	}
	return nil, nil
}

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
	goodsSvc := app.MustGetGoodsService()
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
							Md5fileStr32, err := resourceSvc.Md5SumFile(ev.Name)
							if err != nil {
								fmt.Println(errors.Wrap(err, "Create resourceSvc.Md5SumFile"))
								continue
							}
							Md5path := Md5fileStr32[:2] + Md5fileStr32 + path.Ext(ev.Name)
							CopyDstfullPath := echoapp.ConfigOpts.Asset.TmpRoot + "/ppt/" + Md5path
							if _, err := copy(ev.Name, CopyDstfullPath); err != nil {
								fmt.Println("copy err")
								continue
							}

							strArr, err := getPPTCoverUrl(echoapp.ConfigOpts.Asset.MyURL + path.Base(ev.Name))
							if err != nil && len(strArr) < 2 {
								fmt.Println(errors.Wrap(err, "Create  getPPTCoverUrl"))
								continue
							}

							data, err := json.Marshal(strArr)
							if err != nil {
								fmt.Println(errors.Wrap(err, "Create json.Marsha"))
								continue
							}
							err = goodsSvc.SaveGoods(&echoapp.Goods{
								Name:       path.Base(ev.Name),
								Price:      0.2,
								GoodType:   path.Ext(ev.Name),
								RealPrice:  0.5,
								Covers:     strArr[0],
								SmallCover: string(data),
								Tags:       path.Dir(ev.Name),
								Pages:      len(strArr),
							})
							if err != nil {
								fmt.Println(errors.Wrap(err, "Create  goodsSvc.SaveGoods"))
								continue
							}
							res, err := goodsSvc.GetGoodsByName(path.Base(ev.Name))
							if err != nil {
								fmt.Println(errors.Wrap(err, "Create  goodsSvc.GetGoodsByName"))
								continue
							}
							err = goodsSvc.SaveTags(&echoapp.Tags{
								GoodsId: res.ID,
								Name:    path.Dir(ev.Name),
							})
							if err != nil {
								fmt.Println(errors.Wrap(err, "Create  goodsSvc.SaveTags"))
								continue
							}
							err = resourceSvc.SaveResource(&echoapp.Resource{

								Name:       path.Base(ev.Name),
								Path:       Md5path,
								Type:       path.Ext(ev.Name),
								Covers:     strArr[0],
								SmallCover: string(data),
								GoodsId:    res.ID,
								Pages:      len(strArr),
							})
							if err != nil {
								fmt.Println(errors.Wrap(err, "Create  resourceSvc.SaveResource"))
								continue
							}

						}
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						fmt.Println("写入文件 : ", ev.Name)
						Md5file, err := resourceSvc.Md5SumFile(ev.Name)
						if err != nil {
							fmt.Println(errors.Wrap(err, "Write  resourceSvc.Md5SumFile"))
							continue
						}
						res, err := resourceSvc.GetResourceByName(path.Base(ev.Name))
						if err != nil {
							fmt.Println(errors.Wrap(err, "Write resourceSvc.GetResourceByName"))
							continue
						}
						res.Path = Md5file[:2] + Md5file + path.Ext(ev.Name)
						err = resourceSvc.ModifyResource(res)
						if err != nil {
							fmt.Println(errors.Wrap(err, "Write resourceSvc.ModifyResource"))
							continue
						}
					}
					// if ev.Op&fsnotify.Remove == fsnotify.Remove {
					// 	fmt.Println("删除文件 : ", ev.Name)
					// 	//如果删除文件是目录，则移除监控
					// 	fi, err := os.Stat(ev.Name)
					// 	if err == nil && fi.IsDir() {
					// 		watch.Remove(ev.Name)
					// 		fmt.Println("删除监控 : ", ev.Name)
					// 	} else {
					// 		res, err := resourceSvc.GetResourceByName(path.Base(ev.Name))
					// 		if err != nil {
					// 			continue
					// 		}
					// 		resourceSvc.DeleteResource(res)
					// 	}
					// }
					if ev.Op&fsnotify.Rename == fsnotify.Rename {
						fmt.Println("重命名文件:", ev.Name)
						res, err := resourceSvc.GetResourceByName(path.Base(ev.Name))
						if err != nil {
							fmt.Println(errors.Wrap(err, "Rename resourceSvc.GetResourceByName"))
							continue
						}
						if err := resourceSvc.DeleteResource(res); err != nil {
							fmt.Println(errors.Wrap(err, "Rename resourceSvc.DeleteResource"))
							continue
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
