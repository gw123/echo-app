IMAGE_TAG_VERSION = 1.5.6
REMOTE_USER_API_TAG = "registry.cn-beijing.aliyuncs.com/gapi/echoapp:$(IMAGE_TAG_VERSION)"
DEFAULT_BUILD_TAG = "1.10.1-alpine"
#DOCKER_BUILD_PATH=/Users/mac/code/docker/images/echoapp
DOCKER_BUILD_PATH=./docker
API_VERSION = v$(IMAGE_TAG_VERSION)

ifeq "$(MODE)" "dev"
	BUILD_TAG = "1.10.1"
endif

ifeq "$(BUILD_TAG)" ""
	BUILD_TAG = $(DEFAULT_BUILD_TAG)
endif

upload-all: upload-user upload-file upload-comment upload-goods upload-order upload-site

upload-qys:
	@scp  config.qys.yaml ubuntu@qys:/data/jobs/config.yaml\ &&\
	 scp  upload ubuntu@qys:/data/jobs

upload-user:
	@scp  config.prod.yaml root@sh2:/data/apps/user/config.yaml\ &&\
	 scp  upload  util/scripes/user.sh root@sh2:/data/apps/user

upload-file: file-dir
	@scp  config.prod.yaml root@sh2:/data/apps/file/config.yaml\ &&\
	 scp  upload util/scripes/file.sh root@sh2:/data/apps/file/ &&\
     scp -r resources/storage/keys/  root@sh2:/data/apps/file/resources/storage

upload-order: order-dir
	@scp  config.prod.yaml root@sh2:/data/apps/order/config.yaml\ &&\
	 scp upload util/scripes/order.sh root@sh2:/data/apps/order/ &&\
     scp -r resources/storage/keys/  root@sh2:/data/apps/order/resources/storage

upload-comment: comment-dir
	@scp  config.prod.yaml root@sh2:/data/apps/comment/config.yaml\ &&\
	 scp upload  util/scripes/comment.sh root@sh2:/data/apps/comment\ &&\
     scp -r resources/storage/keys/  root@sh2:/data/apps/comment/resources/storage

upload-site: site-dir
	@ssh root@sh2 cp /data/apps/site/config.yaml /data/apps/site/config.yaml.back\ &&\
	 scp config.prod.yaml root@sh2:/data/apps/site/config.yaml\ &&\
	 scp upload  util/scripes/site.sh root@sh2:/data/apps/site\ &&\
     scp -r resources/public root@sh2:/data/apps/site/resources/ &&\
     scp -r resources/storage/keys/  root@sh2:/data/apps/site/resources/storage\ &&\
     scp -r resources/views/ root@sh2:/data/apps/site/resources/views

upload-goods: goods-dir
	@scp  config.prod.yaml root@sh2:/data/apps/goods/config.yaml\ &&\
	 scp upload  util/scripes/goods.sh root@sh2:/data/apps/goods\ &&\
     scp -r resources/storage/keys/  root@sh2:/data/apps/goods/resources/storage

upload-activity: activity-dir
	@scp  config.prod.yaml root@sh2:/data/apps/activity/config.yaml\ &&\
	 scp upload  util/scripes/activity.sh root@sh2:/data/apps/activity\ &&\
     scp -r resources/storage/keys/  root@sh2:/data/apps/activity/resources/storage

restart:
	ssh root@sh2 supervisorctl reload

upload-user-config:
	scp config.prod.yaml root@sh2:/data/apps/user/

file-dir:
	ssh root@sh2 mkdir -p /data/apps/file
	ssh root@sh2 mkdir -p /data/apps/file/resources/storage

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags  '-w -s' -o echoapp ./entry/main.go &&\
	upx -9 -f -o upload ./echoapp

build-alpine:
	@docker run --rm -v "$(PWD)":/go/src/github.com/gw123/echo-app \
	    -e GOPROXY=https://goproxy.cn \
	    -e GOPRIVATE=github.com/gw123/echo-app \
		-w /go/src/github.com/gw123/echo-app \
		golang:1.15.2-alpine3.12 \
		go build -v -ldflags '-w -s' -o echoapp github.com/gw123/echo-app/entry &&\
		upx -6 -f -o upload ./echoapp

build-static:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags  '-w -s  -extldflags "-static"' -o echoapp ./entry/main.go &&\
	upx -6 -f -o upload ./echoapp

docker-image:
	mkdir -p $(DOCKER_BUILD_PATH)/resources/views &&\
    mkdir -p $(DOCKER_BUILD_PATH)/resources/storage &&\
	cp -r resources/views/ $(DOCKER_BUILD_PATH)/resources/views &&\
    chmod +x upload &&\
    cp upload $(DOCKER_BUILD_PATH)/echoapp &&\
    cp Dockerfile $(DOCKER_BUILD_PATH)/ &&\
    cp -r resources/public/ $(DOCKER_BUILD_PATH)/resources/public &&\
    cp -r resources/storage/keys/ $(DOCKER_BUILD_PATH)/resources/storage
	@docker build -t $(REMOTE_USER_API_TAG) $(DOCKER_BUILD_PATH)
	@docker push $(REMOTE_USER_API_TAG)

## 在宿主机器上静态打包， 打包体积大但是速度快， 适合开发阶段
docker-all: build-static docker-image set-config

## 借助docker容器打包，打包体积小但是速度慢 ，适合正式环境使用
docker-prod: build-alpine docker-image

docker-compose-up:
	cd docker && docker-compose down &&\
	export ECHOAPP_TAG=$(IMAGE_TAG_VERSION) && docker-compose up

run-docker:
	docker run -it --rm  -v $(DOCKER_BUILD_PATH)/etc:/etc/echoapp \
    -v $(DOCKER_BUILD_PATH)/resources/storage/keys:/usr/local/var/echoapp/resources/storage/keys \
    $(REMOTE_USER_API_TAG)  echoapp file --config  /etc/echoapp/config.prod.yaml

supervisor:
	supervisord -c supervisord.conf

set-config:
	@sed 's@{API_VERSION}@$(API_VERSION)@' ./docker/conf/config.tpl.yaml > ./docker/conf/config.yaml &&\
	cat docker/conf/config.yaml | etcdctl $(AUTH) $(ENDPOINTS) put /xyt/config.prod.yaml

set-dev-config:
	cat config.yaml | etcdctl $(AUTH) $(ENDPOINTS) put /xyt/config.yaml

goose:
	goose -dir migrations mysql  'root:password@tcp(sh2:3306)/laraveltest' up

.PHONY: all
all:
	build