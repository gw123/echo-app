package services

import (
	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
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
			name: "001",
			fields: fields{
				ConsumeNoticeUrl: "",
				tongchengOptionMap: map[string]echoapp.TongchengOption{
					"14": {
						Key:    "4MTU1KBG",
						UserId: "9a86097b-b95d-4fd4-bbb9-a18aaafc84b1",
					},
				},
			},
			args: args{
				key: "4MTU1KBG",
				request: echoapp.TongchengRequest{
					RequestHead: echoapp.TongchengRequestHead{
						Sign:      "",
						UserId:    "9a86097b-b95d-4fd4-bbb9-a18aaafc84b1",
						Method:    "ConsumeNotice",
						Version:   "v1.0",
						Timestamp: 1588219228,
					},
					RawRequestBody: `{"orderSerialId":"123","partnerOrderId":"123","consumeDate":"2020-04-08 00:00:00","tickets":1}`,
					//EncryptRequestBody: `sBE4yQDodGqnKpe0BfeLzxdb6ntDQdaRlIEbgsS8OViJwcbTMydj8WVpT8Hgrd3Jq+lT4dz1ULPPWnew344FmLysGcYYLLRF5k1xLiNMaFsJv3ykoK1hao1OuBKZekWp`,
				},
			},
			want: "abf78614e9878708adfd21cee13c7ab3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mSvr := TongchengService{
				ConsumeNoticeUrl:   tt.fields.ConsumeNoticeUrl,
				tongchengOptionMap: tt.fields.tongchengOptionMap,
				mu:                 tt.fields.mu,
			}
			tt.args.request.EncryptRequestBody, _ = echoapp_util.EntryptDesECB([]byte( tt.args.request.RawRequestBody), []byte(tt.args.key))

			if got := mSvr.Sign(tt.args.key, tt.args.request); got != tt.want {
				t.Errorf("Sign() = %v, want %v", got, tt.want)
			}
		})
	}
}
