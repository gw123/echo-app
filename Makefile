DEFAULT_TAG = "echo-app:latest"
DEFAULT_BUILD_TAG = "1.10.1-alpine"

ifeq "$(MODE)" "dev"
	BUILD_TAG = "1.10.1"
endif

ifeq "$(BUILD_TAG)" ""
	BUILD_TAG = $(DEFAULT_BUILD_TAG)
endif

build:
	@docker run --rm -v "$(PWD)":/go/src/github.com/gw123/echo-app \
		-w /go/src/github.com/gw123/echo-app \
		golang:$(BUILD_TAG) \
		go build -v -o echoapp github.com/gw123/echo-app/entry

docker: build
	@docker build -t $(DEFAULT_TAG) .

runUserServer:
	go run entry/main.go user

runGoodsServer:
	go run entry/main.go goods

runOrderServer:
	go run entry/main.go user

.PHONY: all
all:
	build
