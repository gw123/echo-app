APPNAME=echoapp
IMAGE_TAG_VERSION = 1.9.25

REMOTE_USER_API_TAG = "registry.cn-beijing.aliyuncs.com/gapi/$(APPNAME):$(IMAGE_TAG_VERSION)"
REMOTE_USER_API_LATEST_TAG = "registry.cn-beijing.aliyuncs.com/gapi/$(APPNAME):latest"
DEFAULT_BUILD_TAG = "1.10.1-alpine"
#DOCKER_BUILD_PATH=/Users/mac/code/docker/images/$(APPNAME)
DOCKER_BUILD_PATH=./deploy/docker
K8S_PATH= $(PWD)/deploy/k8s/
API_VERSION = v$(IMAGE_TAG_VERSION)

ifeq "$(MODE)" "dev"
	BUILD_TAG = "1.10.1"
endif

ifeq "$(BUILD_TAG)" ""
	BUILD_TAG = $(DEFAULT_BUILD_TAG)
endif

USERNAME=gw123
#USERNAME=root

DEV=dev
SH2=sh2
QYS=qys

HOST=$(DEV)
#HOST=$(SH2)

upload-site: dir
	@ssh $(USERNAME)@$(HOST) cp /data/apps/site/config.yaml /data/apps/site/config.yaml.back\ &&\
	 scp config.prod.yaml root@sh2:/data/apps/site/config.yaml\ &&\
	 scp upload  util/scripes/site.sh root@sh2:/data/apps/site\ &&\
     scp -r resources/public root@sh2:/data/apps/site/resources/ &&\
     scp -r resources/storage/keys/  root@sh2:/data/apps/site/resources/storage\ &&\
     scp -r resources/views/ root@sh2:/data/apps/site/resources/views

upload-view:
	scp -r resources/views/ $(USERNAME)@$(HOST):/data/apps/$(APPNAME)/resources/views

restart:
	ssh $(USERNAME)@$(HOST) supervisorctl reload

dir:
	ssh $(USERNAME)@$(HOST) mkdir -p /data/apps/$(APPNAME)
	ssh $(USERNAME)@$(HOST) mkdir -p /data/apps/$(APPNAME)/resources/storage

init: dir
	scp -r docker/* $(USERNAME)@$(HOST):/data/apps/$(APPNAME)

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags  '-w -s' -o $(APPNAME) ./entry/main.go &&\
	upx -9 -f -o upload ./$(APPNAME)

build-alpine:
	@docker run --rm -v "$(PWD)":/go/src/github.com/gw123/echo-app \
	    -e GOPROXY=https://goproxy.cn \
	    -e GOPRIVATE=github.com/gw123/echo-app \
		-w /go/src/github.com/gw123/echo-app \
		golang:1.15.2-alpine3.12 \
		go build -v -ldflags '-w -s' -o $(APPNAME) github.com/gw123/echo-app/entry &&\
		upx -6 -f -o upload ./$(APPNAME)

build-static:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags  '-w -s  -extldflags "-static"' -o $(APPNAME) ./entry/main.go &&\
	upx -6 -f -o upload ./$(APPNAME)

docker-image:
	mkdir -p $(DOCKER_BUILD_PATH)/resources/views &&\
    mkdir -p $(DOCKER_BUILD_PATH)/resources/storage &&\
	cp -r resources/views/ $(DOCKER_BUILD_PATH)/resources/views &&\
    chmod +x upload &&\
    cp upload $(DOCKER_BUILD_PATH)/$(APPNAME) &&\
    cp Dockerfile $(DOCKER_BUILD_PATH)/ &&\
    cp -r resources/public/ $(DOCKER_BUILD_PATH)/resources/public &&\
    cp -r resources/storage/keys/ $(DOCKER_BUILD_PATH)/resources/storage
	@docker build -t $(REMOTE_USER_API_TAG) $(DOCKER_BUILD_PATH)
	@docker tag $(REMOTE_USER_API_TAG) $(REMOTE_USER_API_LATEST_TAG)
	@docker push $(REMOTE_USER_API_TAG)
	@docker push $(REMOTE_USER_API_LATEST_TAG)

## 在宿主机器上静态打包， 打包体积大但是速度快， 适合开发阶段
docker-all: build-static docker-image

## 借助docker容器打包，打包体积小但是速度慢 ，适合正式环境使用
docker-prod: build-alpine docker-image

docker-compose-up:
	cd docker && docker-compose down &&\
	export ECHOAPP_TAG=$(IMAGE_TAG_VERSION) && docker-compose up

docker-compose-up-remote:
	ssh $(USERNAME)@$(HOST) \
	'cd /data/apps/$(APPNAME) && docker-compose down &&\
	export ECHOAPP_TAG=$(IMAGE_TAG_VERSION) && docker-compose up'

run-docker:
	docker run -it --rm  -v $(DOCKER_BUILD_PATH)/etc:/etc/$(APPNAME) \
    -v $(DOCKER_BUILD_PATH)/resources/storage/keys:/usr/local/var/$(APPNAME)/resources/storage/keys \
    $(REMOTE_USER_API_TAG)  $(APPNAME) file --config  /etc/$(APPNAME)/config.prod.yaml

supervisor:
	supervisord -c supervisord.conf

set-config:
	@sed 's@{API_VERSION}@$(API_VERSION)@' ./docker/conf/config.tpl.yaml > ./docker/conf/config.yaml &&\
	cat docker/conf/config.yaml | etcdctl $(AUTH) $(ENDPOINTS) put /xyt/config.prod.yaml

set-dev-config:
	cat config.yaml | etcdctl $(AUTH) $(ENDPOINTS) put /xyt/config.yaml

goose:
	goose -dir migrations mysql  'root:password@tcp(sh2:3306)/laraveltest' up

update-k8s:
	kubectl replace --force -f $(K8S_PATH)site-svr.yaml&&\
	kubectl replace --force -f $(K8S_PATH)goods-svr.yaml&&\
	kubectl replace --force -f $(K8S_PATH)order-svr.yaml&&\
	kubectl replace --force -f $(K8S_PATH)user-svr.yaml&&\
	kubectl replace --force -f $(K8S_PATH)comment-svr.yaml

tail-k8s:
	#kubectl get pods -nechoapp |grep site |awk  '{print $1}' | kubectl logs -f -nechoapp
	kubectl get pods -nechoapp |grep site |awk   '{print $$1}'

.PHONY: all
all: docker-all docker-compose-up-remote

.PHONY: k8s-all
k8s-all: docker-all  update-k8s

restart_dev: docker-compose-up-remote
