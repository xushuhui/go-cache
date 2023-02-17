package util

import "testing"

func TestParseBytes(t *testing.T) {
	t.Log(1 << 5)
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{

		{"kb", args{s: "1kb"}, 1024},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := ParseBytes(tt.args.s)

			if got != tt.want {
				t.Errorf("ParseBytes() got = %v, want %v", got, tt.want)
			}
		})
	}
}
