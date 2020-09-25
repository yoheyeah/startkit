package types

import (
	"log"
	"reflect"
	"testing"
)

func TestNewVarAddrByInterface(t *testing.T) {
	type a struct {
		value int
	}
	type args struct {
		addr interface{}
	}
	aAddr := a{
		value: 1,
	}
	tests := []struct {
		name string
		args args
		want reflect.Type
	}{
		{
			name: "test get address func",
			args: args{
				addr: &aAddr,
			},
			want: reflect.TypeOf(&aAddr),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewVarAddrByInterface(tt.args.addr)
			if !reflect.DeepEqual(reflect.TypeOf(got), tt.want) {
				t.Errorf("NewVarAddrByInterface() = %v, want %v", got, tt.want)
				log.Println((*(got.(*interface{}))).(a).value)
			} else {
				log.Println((got).(a).value)
			}
		})
	}
}

func TestNewVarByInterface(t *testing.T) {
	type a struct {
		value int
	}
	type args struct {
		obj interface{}
	}
	aAddr := a{
		value: 1,
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test get address func",
			args: args{
				obj: aAddr,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewVarByInterface(tt.args.obj); !reflect.DeepEqual(got.(a).value, tt.want) {
				t.Errorf("NewVarByInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}
