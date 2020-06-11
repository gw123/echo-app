package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/app"
	"github.com/pkg/errors"

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

	resourceSvc := app.MustGetResourceService()
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
							md5fileStr32, err := echoapp_util.Md5SumFile(ev.Name)
							if err != nil {
								fmt.Println(errors.Wrap(err, "Md5SumFile"))
								continue
							}
							md5path := md5fileStr32[:2] + md5fileStr32 + path.Ext(ev.Name)
							if _, err := resourceSvc.GetResourceByMd5Path(nil, md5path); err == nil {
								fmt.Println(errors.New(ev.Name + " :It already has the same content"))
								continue
							}
							fileType := echoapp_util.GetFileType(ev.Name)
							// CopyDstfullPath := echoapp.ConfigOpts.Asset.StorageRoot + "/" + fileType + "/" + md5path
							// if err := echoapp_util.Copy(CopyDstfullPath, ev.Name); err != nil {
							// 	fmt.Println("copy err")
							// 	continue
							// }
							if _, err := echoapp_util.UploadFileToQiniu(ev.Name, "/"+fileType+"/"+md5path); err != nil {
								fmt.Println(errors.Wrap(err, "uploadFileToQiniu"))
								continue
							}
							strArr, err := echoapp_util.GetPPTCoverUrl(echoapp.ConfigOpts.ResourceOptions.BaseURL + "/" + fileType + "/" + path.Base(ev.Name))
							if err != nil && len(strArr) < 2 {
								fmt.Println(errors.Wrap(err, "Create  getPPTCoverUrl"))
								continue
							}

							data, err := json.Marshal(strArr)
							if err != nil {
								fmt.Println(errors.Wrap(err, "Create json.Marsha"))
								continue
							}
							err = goodsSvc.Save(&echoapp.Goods{
								GoodsBrief: echoapp.GoodsBrief{
									Name:  path.Base(ev.Name),
									Price: 0.2,
									//GoodType:   path.Ext(ev.Name),
									RealPrice:  0.5,
									Covers:     strArr[0],
									SmallCover: string(data),
									Tags:       path.Dir(ev.Name),
									//Pages:      len(strArr),
								},
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
							err = goodsSvc.SaveTag(&echoapp.GoodsTag{
								//GoodsId: res.ID,
								Name: path.Dir(ev.Name),
							})
							if err != nil {
								fmt.Println(errors.Wrap(err, "Create  goodsSvc.SaveTags"))
								continue
							}

							res_tag, err := goodsSvc.GetTagByName(path.Base(ev.Name))
							if err != nil {
								fmt.Println(errors.Wrap(err, "Create goodsSvc.GetTagsByName"))
								continue
							}
							err = resourceSvc.SaveResource(&echoapp.Resource{
								TagId:      int64(res_tag.ID),
								Name:       path.Base(ev.Name),
								Path:       md5path,
								Type:       path.Ext(ev.Name),
								Covers:     strArr[0],
								SmallCover: string(data),
								GoodsId:    int64(res.ID),
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
						Md5file, err := echoapp_util.Md5SumFile(ev.Name)
						if err != nil {
							fmt.Println(errors.Wrap(err, "Write  Md5SumFile"))
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

	echoapp_util.DefaultLogger().Infof("开始监控")

	go watchDir(echoapp.ConfigOpts.Asset.StorageRoot)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

}

var watcherCmd = &cobra.Command{
	Use:   "watch",
	Short: "监控",
	Long:  "监听目录",
	Run: func(cmd *cobra.Command, args []string) {
		startWatcher()
	},
}

func init() {
	rootCmd.AddCommand(watcherCmd)
}
