package main

import (
	"encoding/json"
	"testing"
	"time"
)

func Test_distributeOrders(t *testing.T) {
	orders := buildMockOrders()
	type args struct {
		orders   []Order
		initData time.Time
		endData  time.Time
	}
	tests := []struct {
		name string
		args args
		want map[string]int
	}{
		{
			name: "Should successfully execute and distribute orders",
			args: args{
				orders:   orders,
				initData: time.Now().Add(monthsToSubtract(30)),
				endData:  time.Now(),
			},
			want: map[string]int{
				Range1and3Months:         30,
				Range4and6Months:         30,
				Range7and12Months:        60,
				RangeGreaterThan12Months: 170,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := distributeOrders(tt.args.orders, tt.args.initData, tt.args.endData)
			want, err := json.Marshal(tt.want)
			if err != nil {
				t.Error(err)
			}
			gotBytes, err := json.Marshal(got)
			if err != nil {
				t.Error(err)
			}

			if string(want) != string(gotBytes) {
				t.Error("mismatch result")
			}
		})
	}
}
