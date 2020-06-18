package echoapp_util

import (
	"reflect"
	"testing"
)

// args: args{
// 	a:      1.0,
// 	b:      2.0,
// 	Ptheta: []float64{0.5, 0.5, 0.3, 0.4},
// 	wight:  []float64{0.25, 0.1, 0.3, 0.35},
// },
// want: 0.39884445629608023,
// },
// {
// args: args{
// 	a:      1.0,
// 	b:      2.0,
// 	Ptheta: []float64{0.7, 0.7, 0.6, 0.6},
// 	wight:  []float64{0.25, 0.1, 0.3, 0.35},
// },
// want: 0.6326838983950979,
// },

func TestLinguisticToTFS(t *testing.T) {
	type args struct {
		value int
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		{
			args: args{
				value: 0,
			},
			want: []float64{0, 0, 0.25},
		},
		{
			args: args{
				value: 1,
			},
			want: []float64{0, 0.25, 0.5},
		},
		{
			args: args{
				value: 3,
			},
			want: []float64{0.5, 0.75, 1},
		},
		{
			args: args{
				value: 5,
			},
			want: []float64{0.5, 0.75, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LinguisticToTFS(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LinguisticToTFS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTFSToFS(t *testing.T) {
	type args struct {
		tfn []float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			args: args{
				tfn: []float64{0.5, 0.75, 1},
			},
			want: 0.75,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TFSToFS(tt.args.tfn); got != tt.want {
				t.Errorf("TFSToFS() = %v, want %v", got, tt.want)
			}
		})
	}
}
