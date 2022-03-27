package main

import "testing"

func Test_calcDif(t *testing.T) {
	type args struct {
		given  float64
		target float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"test1", args{given: 1, target: 2}, 0.5},
		{"test2", args{given: 40, target: 2}, 0.05},
		{"test3", args{given: 10, target: 2}, 0.2},
		{"test4", args{given: 0.5, target: 2}, 0.25},
		{"test5", args{given: 2, target: 2}, 1},
		{"test6", args{given: 1, target: 1}, 1},
		{"test7", args{given: 3, target: 2}, 0.6666666666666666}, //or this might be 0.5?
		{"test8", args{given: 2, target: 3}, 0.6666666666666666},
		{"test8", args{given: 1000, target: 3}, 0.003},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcDif(tt.args.given, tt.args.target); got != tt.want {
				t.Errorf("calcDif() = %v, want %v", got, tt.want)
			}
		})
	}
}
