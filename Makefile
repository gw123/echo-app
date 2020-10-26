DEFAULT_TAG = "echo-app:1.0.2"
REMOTE_USER_API_TAG = "registry.cn-beijing.aliyuncs.com/gapi/user:1.0.1"
DEFAULT_BUILD_TAG = "1.10.1-alpine"
DOCKER_BUILD_PATH=/data/docker/images/echoapp

ifeq "$(MODE)" "dev"
	BUILD_TAG = "1.10.1"
endif

ifeq "$(BUILD_TAG)" ""
	BUILD_TAG = $(DEFAULT_BUILD_TAG)
endif

<<<<<<< HEAD
upload-user:
	@scp -r resources/views/ root@sh2:/data/apps/user/resources/views &&\
     scp  config.yaml upload util/scripes/user.sh root@sh2:/data/apps/user/ &&\
     scp -r resources/public root@sh2:/data/apps/user/resources/public &&\
     scp -r resources/storage/keys/  root@sh2:/data/apps/user/resources/storage

upload-file: file-dir
	@scp  config.yaml upload util/scripes/file.sh root@sh2:/data/apps/file/ &&\
     scp -r resources/public  root@sh2:/data/apps/file/resources/public &&\
     scp -r resources/storage/keys/  root@sh2:/data/apps/file/resources/storage

upload-order: order-dir
	@scp  config.yaml upload util/scripes/order.sh root@sh2:/data/apps/order/ &&\
     scp -r resources/public  root@sh2:/data/apps/order/resources/public &&\
     scp -r resources/storage/keys/  root@sh2:/data/apps/order/resources/storage

upload-comment: comment-dir
	@scp  config.yaml upload  util/scripes/comment.sh root@sh2:/data/apps/comment\ &&\
     scp -r resources/public  root@sh2:/data/apps/comment/resources/public &&\
     scp -r resources/storage/keys/  root@sh2:/data/apps/comment/resources/storage

upload-site: site-dir
	@scp  config.yaml upload  util/scripes/site.sh root@sh2:/data/apps/site\ &&\
     scp -r resources/public  root@sh2:/data/apps/site/resources/ &&\
     scp -r resources/storage/keys/  root@sh2:/data/apps/site/resources/storage

restart:
	ssh root@sh2 supervisorctl restart user
=======
upload-all: upload-user upload-file upload-comment upload-goods upload-order upload-site

upload-user:
	@scp  config.prod.yaml root@sh2:/data/apps/user/config.yaml\ &&\
	 scp  upload  util/scripes/user.sh root@sh2:/data/apps/user\ &&\
     scp  upload util/scripes/user.sh root@sh2:/data/apps/user/

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
	@scp  config.prod.yaml root@sh2:/data/apps/site/config.yaml\ &&\
     scp -r resources/public root@sh2:/data/apps/site/resources/public &&\
	 scp upload  util/scripes/site.sh root@sh2:/data/apps/site\ &&\
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
	scp  config.prod.yaml root@sh2:/data/apps/user/
>>>>>>> 9fc3ae585f42ae0efeb5f967a4ff431a7b4509d4

file-dir:
	ssh root@sh2 mkdir -p /data/apps/file
	ssh root@sh2 mkdir -p /data/apps/file/resources/storage

<<<<<<< HEAD
=======
activity-dir:
	ssh root@sh2 mkdir -p /data/apps/activity
	ssh root@sh2 mkdir -p /data/apps/activity/resources/storage

>>>>>>> 9fc3ae585f42ae0efeb5f967a4ff431a7b4509d4
user-dir:
	ssh root@sh2 mkdir -p /data/apps/user
	ssh root@sh2 mkdir -p /data/apps/user/resources/storage

comment-dir:
	ssh root@sh2 mkdir -p /data/apps/comment
	ssh root@sh2 mkdir -p /data/apps/comment/resources/storage

order-dir:
	ssh root@sh2 mkdir -p /data/apps/order
	ssh root@sh2 mkdir -p /data/apps/order/resources/storage

site-dir:
	ssh root@sh2 mkdir -p /data/apps/site
	ssh root@sh2 mkdir -p /data/apps/site/resources/storage

<<<<<<< HEAD
build:
	go build -ldflags  '-w -s' -o echoapp ./entry/main.go &&\
	rm upload && upx -9 -o upload ./echoapp
=======
goods-dir:
	ssh root@sh2 mkdir -p /data/apps/goods
	ssh root@sh2 mkdir -p /data/apps/goods/resources/storage
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags  '-w -s' -o echoapp ./entry/main.go &&\
	upx -9 -f -o upload ./echoapp
>>>>>>> 9fc3ae585f42ae0efeb5f967a4ff431a7b4509d4

build-alpine:
	@docker run --rm -v "$(PWD)":/go/src/github.com/gw123/echo-app \
	    -e GOPROXY=https://goproxy.cn \
	    -e GOPRIVATE=github.com/gw123/echo-app \
		-w /go/src/github.com/gw123/echo-app \
		golang:$(BUILD_TAG) \
		go build -v -ldflags '-w -s' -o echoapp github.com/gw123/echo-app/entry &&\
		rm upload && upx -9 -o upload ./echoapp

docker: build-alpine
	mkdir -p $(DOCKER_BUILD_PATH)/resources/views &&\
	mkdir -p $(DOCKER_BUILD_PATH)/etc &&\
    mkdir -p $(DOCKER_BUILD_PATH)/resources/storage &&\
	cp -r resources/views/ $(DOCKER_BUILD_PATH)/resources/views &&\
    cp upload $(DOCKER_BUILD_PATH)/echoapp &&\
    cp Dockerfile $(DOCKER_BUILD_PATH)/ &&\
<<<<<<< HEAD
    cp config.docker.yaml $(DOCKER_BUILD_PATH)/etc/config.yaml &&\
=======
    cp config.docker.yaml $(DOCKER_BUILD_PATH)/etc/config.prod.yaml &&\
>>>>>>> 9fc3ae585f42ae0efeb5f967a4ff431a7b4509d4
    cp -r resources/public/ $(DOCKER_BUILD_PATH)/resources/ &&\
    cp -r resources/storage/keys/ $(DOCKER_BUILD_PATH)/resources/storage
	@docker build -t $(REMOTE_USER_API_TAG) $(DOCKER_BUILD_PATH)
	@docker push $(REMOTE_USER_API_TAG)
	@rm -f echoapp

run-docker:
	docker run -it --rm  -v $(DOCKER_BUILD_PATH)/etc:/etc/echoapp \
    -v $(DOCKER_BUILD_PATH)/resources/storage/keys:/usr/local/var/echoapp/resources/storage/keys \
<<<<<<< HEAD
    $(REMOTE_USER_API_TAG)  echoapp file --config  /etc/echoapp/config.yaml
=======
    $(REMOTE_USER_API_TAG)  echoapp file --config  /etc/echoapp/config.prod.yaml
>>>>>>> 9fc3ae585f42ae0efeb5f967a4ff431a7b4509d4

runUserServer:
	go run entry/main.go user

runGoodsServer:
	go run entry/main.go goods

runOrderServer:
	go run entry/main.go user

supervisor:
	supervisord -c supervisord.conf

goose:
	goose -dir migrations mysql  'root:password@tcp(sh2:13306)/laraveltest' up

.PHONY: all
all:
	build
