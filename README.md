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

## jws中间件认证
    配置文件:
    ```
    jws:
      audience: "user"
      issuer: "xytschool"
      timeout: 360000
      public_key_path: "./resources/keys/jws_public_key.pem"
      private_key_path: "./resources/keys/jws_private_key.pem"
      hash_ids_salt: "123456"
    ```
    注意私钥只有在生成jws签名的模块或者服务需要加载
    只做签名认证的服务或者只需加载公钥就可以,并且私钥的路径需要配置为空.  
	jwsMiddleware := echoapp_middlewares.NewJwsMiddlewares(middleware.DefaultSkipper, app.MustGetJwsHelper())
	jwsAuth.Use(jwsMiddleware)
	
## jws认证后可以使用 echoapp_util.GetCtxtUserId(ctx)获取用户id,想要获取用户完整信息需要配合user中间件,但是注意user中间件所在的服务需要有访问user数据库的权限,在一些微服务场景下我们往往只关心userId

## user中间件获取用户信息, 在配置中间件的action中就可以使用echoapp_util.GetCtxtUser(ctx)获取用户信息  
	userMiddleware := echoapp_middlewares.NewUserMiddlewares(middleware.DefaultSkipper, usrSvr)
	authgoup.Use(userMiddleware)
	
## 获取用户信息 echoapp_util.GetCtxtUser(ctx)
```
   func (sCtl *UserController) GetUserInfo(ctx echo.Context) error {
   	echoapp_util.ExtractEntry(ctx).Info("getUserInfo")
   	user, err := echoapp_util.GetCtxtUser(ctx)
   	if err != nil {
   		return sCtl.Fail(ctx, echoapp.Err_NotFound, "未发现用户", err)
   	}
   	return sCtl.Success(ctx, user)
   }

```