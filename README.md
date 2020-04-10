#基于echo框架的web应用模板

#编译使用
- go mod vendor
- go run entry/main.go

## app/EchoApp 是一个常用接口集合, EchoApp所有属性都是私有防止未初始化被使用，
## 正确的使用方式是调用 GetXXX 和MustGetXXX这样的函数获取接口的实例对象,
## 在GETXXX函数里面使用单例模式方式重复创建

## 记录日志方法，通过下面方式记录可以把 request_id 记录到日志中方便同一个请求链路追踪, 并且还可以导出其他ctx中记录的信息
echoapp_util.ExtractEntry(ctx).Info(renderParams)


