package echoapp_util

import "testing"

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EntryptDesECB(tt.args.data, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("EntryptDesECB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EntryptDesECB() got = %v, want %v", got, tt.want)
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
