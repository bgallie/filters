package ascii85

import (
	"io"
	"strings"
	"testing"
)

func TestToASCII85(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestOne",
			args: args{r: strings.NewReader("This is only a test")},
			want: "<+oue+DGm>Df0B:+CQC7ATMq",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := io.ReadAll(ToASCII85(tt.args.r)); strings.Compare(string(got), tt.want) != 0 {
				t.Errorf("ToASCII85() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestFromASCII85(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestOne",
			args: args{r: strings.NewReader("<+oue+DGm>Df0B:+CQC7ATMq")},
			want: "This is only a test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := io.ReadAll(FromASCII85(tt.args.r)); strings.Compare(string(got), tt.want) != 0 {
				t.Errorf("FromASCII85() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
