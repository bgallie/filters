package hex

import (
	"io"
	"strings"
	"testing"
)

func TestToHex(t *testing.T) {
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
			args: args{r: strings.NewReader("This is only a test of the ToHex filter.")},
			want: "54686973206973206f6e6c7920612074657374206f662074686520546f4865782066696c7465722e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := io.ReadAll(ToHex(tt.args.r)); strings.Compare(string(got), tt.want) != 0 {
				t.Errorf("ToHex() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestFromHex(t *testing.T) {
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
			args: args{r: strings.NewReader("54686973206973206f6e6c7920612074657374206f66207468652046726f6d4865782066696c7465722e")},
			want: "This is only a test of the FromHex filter.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := io.ReadAll(FromHex(tt.args.r)); strings.Compare(string(got), tt.want) != 0 {
				t.Errorf("FromHex() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
