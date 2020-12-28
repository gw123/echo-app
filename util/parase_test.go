package echoapp_util

import (
	"reflect"
	"testing"
)

func TestFetchImageUrls(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "001",
			args: args{
				text: `<div class="container">
            <div class="img-area">
            <img class="my-photo" alt="loading" title="查看大图" style="cursor: pointer;" data-src="http://help.xytschool.com:8012/test/0.jpg" src="images/loading.gif" onclick="changePreviewType('allImages')">
        </div>
        <div class="img-area">
            <img class="my-photo" alt="loading" title="查看大图" style="cursor: pointer;" data-src="http://help.xytschool.com:8012/test/1.jpg" src="images/loading.gif" onclick="changePreviewType('allImages')">
        </div></div>`,
			},
			want: []string{
				"http://help.xytschool.com:8012/test/0.jpg",
				"http://help.xytschool.com:8012/test/1.jpg",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FetchImageUrls(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchImageUrls() = %v, want %v", got, tt.want)
			}
		})
	}
}
