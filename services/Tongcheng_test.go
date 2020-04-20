package services

import (
	echoapp "github.com/gw123/echo-app"
	"sync"
	"testing"
)

func TestTongchengService_Sign(t *testing.T) {
	type fields struct {
		ConsumeNoticeUrl   string
		tongchengOptionMap map[string]echoapp.TongchengOption
		mu                 sync.Mutex
	}
	type args struct {
		key     string
		request echoapp.TongchengRequest
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "001",
			fields: fields{
				ConsumeNoticeUrl:   "",
				tongchengOptionMap: map[string]echoapp.TongchengOption{
					"14" : {
						Key:    "4MTU1KBG",
						UserId: "be3ItcG2WLeCbZVDSjIQG6p6ygFOGr",
					},
				},
			},
			args:   args{
				key:     "",
				request: echoapp.TongchengRequest{},
			},
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mSvr := TongchengService{
				ConsumeNoticeUrl:   tt.fields.ConsumeNoticeUrl,
				tongchengOptionMap: tt.fields.tongchengOptionMap,
				mu:                 tt.fields.mu,
			}
			if got := mSvr.Sign(tt.args.key, tt.args.request); got != tt.want {
				t.Errorf("Sign() = %v, want %v", got, tt.want)
			}
		})
	}
}