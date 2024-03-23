package main

import (
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

//my goal is to become a gopher
func myGoal(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RequestURI)
	spew.Dump(*r.URL)
	_, _ = w.Write([]byte("I wan`t to become a gopher."))
}

func main() {

	//1.注册一个处理器函数
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", myGoal)

	//2.设置监听的TCP地址并启动服务
	//参数1:TCP地址(IP+Port)
	//参数2:handler 创建新的*serveMux,不使用默认的
	err := http.ListenAndServe("127.0.0.1:9009", serveMux)
	if err != nil {
		fmt.Printf("http.ListenAndServe()函数执行错误,错误为:%v\n", err)
		return
	}
}
