package echoapp_util

import (
	"reflect"
	"testing"

	"github.com/labstack/echo"
)

func TestGetFileType(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "./res/hg.ppt",
			args: args{
				filename: "./res/hg.ppt",
			},
			want: "ppt",
		},
		{
			name: "001",

			args: args{
				filename: "./resources/tmp/ppt/cp.ppt",
			},
			want: "ppt",
		},
		{
			name: "002",

			args: args{
				filename: "./resources/tmp/ppt/ch.ppt",
			},
			want: "ppt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFileType(tt.args.filename); got != tt.want {
				t.Errorf("GetFileType() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestDoHttpRequest(t *testing.T) {
// 	type args struct {
// 		url    string
// 		method string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    []byte
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := DoHttpRequest(tt.args.url, tt.args.method)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("DoHttpRequest() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("DoHttpRequest() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestMd5SumFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "001",

			args: args{
				filename: "../resources/tmp/ppt/cp.ppt",
			},
			want: "1e9948e3add3efe730eebb79e7e1b26f",
			//wantErr: ,
		},
		{
			name: "002",

			args: args{
				filename: "../resources/tmp/ppt/ch.ppt",
			},
			want: "77b3a6ca2359a967cdf7217bc3f7bb68",
			//wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Md5SumFile(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("Md5SumFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Md5SumFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUploadFile(t *testing.T) {
	type args struct {
		c           echo.Context
		formname    string
		uploadpath  string
		maxfilesize int64
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UploadFile(tt.args.c, tt.args.formname, tt.args.uploadpath, tt.args.maxfilesize)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UploadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUploadFileToQiniu(t *testing.T) {
	type args struct {
		localFile string
		key       string
	}
	tests := []struct {
		name    string
		args    args
		want    *MyPutRet
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UploadFileToQiniu(tt.args.localFile, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadFileToQiniu() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UploadFileToQiniu() = %v, want %v", got, tt.want)
			}
		})
	}
}
