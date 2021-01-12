#!/usr/bin/env sh

# Install proto3 from source macOS only.
#  brew install autoconf automake libtool
#  git clone https://github.com/google/protobuf
#  ./autogen.sh ; ./configure ; make ; make install
#
# Update protoc Go bindings via
#  go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
#
# See also
#  https://github.com/grpc/grpc-go/tree/master/examples

protoc *.proto --proto_path=. --proto_path=$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis --go_out=plugins=grpc:. --grpc-gateway_out=logtostderr=true:.
protoc gather.proto --proto_path=. --proto_path=$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis --swagger_out=logtostderr=true:../
ls *.pb.go | xargs -n1 -IX bash -c 'protoc-go-inject-tag -input=X'
ls *.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'
