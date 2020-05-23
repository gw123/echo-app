package echoapp_util

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/labstack/gommon/log"
	"testing"
)

func TestEntryptDesECB(t *testing.T) {
	type args struct {
		data []byte
		key  []byte
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
				data: []byte(`{"orderSerialId":"231231","partnerOrderId":"12312312","consumeDate":"2020-04-16 00:00:00","tickets":1}`),
				key:  []byte("4MTU1KBG"),
			},
			want:    "sBE4yQDodGqnKpe0BfeLz6dr/pxczd9O2G0BP0Wr3h5vUGEYgjsROopUNYs8OLRzMo4jYO9SdAWr6VPh3PVQs89ad7DfjgWYl932GXj2y/nmTXEuI0xoWwm/fKSgrWFqjU64Epl6Rak=",
			wantErr: false,
		},
		{
			name: "002",
			args: args{
				data: []byte(`{"orderSerialId":"sz5ua7ffud21255bu85498942","partnerOrderId":"sz5ua7ffud21255bu85498942","consumeDate":"2020-04-29 23:59:08","tickets":1}`),
				key:  []byte("4MTU1KBG"),
			},
			want:    "sBE4yQDodGqnKpe0BfeLz9r/Z/m+DaB/ZQ8wrbzCRxdUO33noZhD6Wq8qU9oOpbpRqoYkaUcJCWNkEFQCsNZThPtQ/JR9CnymJ1pDlQ8kybYVKI69mY9FSpHxTNdO07yKfgJMJ76ssBA7dCzsvLPfEduzxBF63n7mNvgUMAZ+YeUAvV5WWpQJ2+Ury7AQbQE",
			wantErr: false,
		},
		{
			name: "003",
			args: args{
				data: []byte(`{"orderSerialId":"123","partnerOrderId":"123","consumeDate":"2020-04-08 00:00:00","tickets":1}`),
				key:  []byte("4MTU1KBG"),
			},
			want:    "sBE4yQDodGqnKpe0BfeLzxdb6ntDQdaRlIEbgsS8OViJwcbTMydj8WVpT8Hgrd3Jq+lT4dz1ULPPWnew344FmLysGcYYLLRF5k1xLiNMaFsJv3ykoK1hao1OuBKZekWp",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EntryptDesECB(tt.args.data, tt.args.key)

			if (err != nil) != tt.wantErr {
				t.Errorf("EntryptDesECB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EntryptDesECB() \ngot = %v\nwant= %v", got, tt.want)
			}
		})
	}
}

func TestDecryptDESECB(t *testing.T) {
	type args struct {
		d   []byte
		key []byte
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
				d:   []byte("sBE4yQDodGqnKpe0BfeLz6dr/pxczd9O2G0BP0Wr3h5vUGEYgjsROopUNYs8OLRzMo4jYO9SdAWr6VPh3PVQs89ad7DfjgWYl932GXj2y/nmTXEuI0xoWwm/fKSgrWFqjU64Epl6Rak="),
				key: []byte("4MTU1KBG"),
			},
			want:    `{"orderSerialId":"231231","partnerOrderId":"12312312","consumeDate":"2020-04-16 00:00:00","tickets":1}`,
			wantErr: false,
		},
		{
			name: "002",
			args: args{
				d:   []byte("8woZ6zkYhB/b+KFhWbRzzyCIzhVMAlwrfmKVUnWxKNE="),
				key: []byte("4MTU1KBG"),
			},
			want:    `{"refundStatus":1,"remark":""}`,
			wantErr: false,
		},
		{
			name: "003",
			args: args{
				d:   []byte("8woZ6zkYhB/b+KFhWbRzzyCIzhVMAlwrfmKVUnWxKNE=1"),
				key: []byte("4MTU1KBG"),
			},
			want:    ``,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecryptDESECB(tt.args.d, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecryptDESECB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DecryptDESECB() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMd5(t *testing.T) {
	body := `{"requestHead":{"sign":"1b2d1e6374c0a01bb4212b50d0088862","user_id":"df83b69a-cafc-4001-bc66-ac458179f0dd","method":"ConsumeNotice","version":"v1.0","timestamp":1588223586},"requestBody":{"tickets":1,"orderSerialId":"sz112153048963144089698985","partnerOrderId":"sz112153048963144089698985","consumeDate":"2020-04-30 13:13:06"}}`
	data := []byte(body)
	md5str1 := fmt.Sprintf("%x", md5.Sum(data))
	log.Info(md5str1)

	hash1 := md5.New()
	hash1.Write(data)
	log.Info(base64.StdEncoding.EncodeToString(hash1.Sum(nil)))

}
