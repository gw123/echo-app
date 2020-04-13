#基于echo框架的web应用模板,本模板一个echo的实用实例，主要目的是为了提供一个构建代码友好可维护性高的echo应用
- 友好而实用的日志记录方式
- 友好service调用方式

#编译使用
- go mod vendor
- go run entry/main.go

## app/EchoApp 是一个常用接口集合, EchoApp所有属性都是私有防止未初始化被使用，
## 在app.go中GetServiceName函数里面使用单例模式方式防止重复创建
## 要想获取servece正确的使用方式是调用GetServiceName和MustGetServiceName这样的函数获取接口的实例对象,
## 使用获取service方式实例 参考controllers/area.go
```go
	areaMap, err := app.MustGetAreaService().GetAreaMap(areaId)
	if err != nil {
		return c.BaseController.Fail(ctx, echoapp.Error_ArgumentError, "", err)
	}
```

## 记录日志方法，通过下面方式记录可以把 request_id 记录到日志中方便同一个请求链路追踪, 并且还可以导出其他ctx中记录的信息
echoapp_util.ExtractEntry(ctx).Info(renderParams)


