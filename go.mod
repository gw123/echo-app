module github.com/gw123/echo-app

go 1.15

require (
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.246
	github.com/aymerick/raymond v2.0.2+incompatible
	github.com/bsm/redislock v0.5.0
	github.com/chanxuehong/wechat v0.0.0-20200409104612-0a1fd76d7a3a
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/davecgh/go-spew v1.1.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/forgoer/openssl v0.0.0-20200331032942-ad9f8d57d8b1
	github.com/fsnotify/fsnotify v1.4.9
	github.com/go-kit/kit v0.10.0
	github.com/go-redis/redis/v7 v7.4.0
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/gw123/glog v1.3.2
	github.com/gw123/gworker v1.6.1
	github.com/iGoogle-ink/gopay v1.5.16-0.20200714134502-68dab747848e
	github.com/iGoogle-ink/gotil v1.0.3
	github.com/jinzhu/gorm v1.9.14
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.3.0
	github.com/lib/pq v1.7.0 // indirect
	github.com/lightstep/lightstep-tracer-go v0.18.1
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/medivhzhan/weapp v1.5.1
	github.com/oklog/oklog v0.3.2
	github.com/olivere/elastic/v7 v7.0.19
	github.com/opentracing/opentracing-go v1.2.0
	github.com/philchia/agollo/v4 v4.1.2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.3.0
	github.com/prometheus/common v0.10.0
	github.com/qiniu/api.v7/v7 v7.4.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/shima-park/agollo v1.2.10
	github.com/silenceper/wechat/v2 v2.0.1
	github.com/sirupsen/logrus v1.7.0
	github.com/skip2/go-qrcode v0.0.0-20200526175731-7ac0b40b2038
	github.com/speps/go-hashids v2.0.0+incompatible
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	github.com/streadway/amqp v0.0.0-20200108173154-1c71cc93ed71
	github.com/zouyx/agollo/v4 v4.0.3
	golang.org/x/net v0.0.0-20201209123823-ac852fbbde11
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/genproto v0.0.0-20201119123407-9b1e624d6bc4
	google.golang.org/grpc v1.33.1
	sourcegraph.com/sourcegraph/appdash v0.0.0-20190731080439-ebfcffb1b5c0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.29.1
