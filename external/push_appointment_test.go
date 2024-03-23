package external

import (
	"context"
	"testing"

	"github.com/gw123/glog"
)

func TestDoPushAppointmentRequest(t *testing.T) {
	ctx := context.Background()

	type args struct {
		request *PushAppointmentRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "01",
			args: args{request: &PushAppointmentRequest{
				ScenicID: "1300000",
				BookList: []map[string][]BookItem{
					{
						"2021-01-09": []BookItem{{
							Label:     "08:00-09:00",
							StartTime: 1590278400,
							EndTime:   1590278400,
							BookNum:   9,
							RemainNum: 22,
						},
						},
					},
				},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DoPushAppointmentRequest(ctx, []*PushAppointmentRequest{tt.args.request})
			if (err != nil) != tt.wantErr {
				t.Errorf("DoPushAppointmentRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			glog.DefaultLogger().Infof("DoPushAppointmentRequest over")
		})
	}
}
