package rdb

import (
	"reflect"
	"testing"
)

func TestInitRdb(t *testing.T) {
	tests := []struct {
		name string
		want *Rdb
	}{
			{"local", &Rdb{
				Host:         "172.19.192.1",
				Port:         "5432",
				User:         "admin1",
				Password:     "123456123",
				DatabaseName: "myapp",
				DriverName:   "postgres",
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitRdbInfo(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitRdbInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
