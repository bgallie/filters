package filters

import (
	"io"
	"strings"
	"testing"
)

func TestToBase64(t *testing.T) {
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
			args: args{r: strings.NewReader("This is only a test of the ToBase64 filter.")},
			want: "VGhpcyBpcyBvbmx5IGEgdGVzdCBvZiB0aGUgVG9CYXNlNjQgZmlsdGVyLg==",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := io.ReadAll(ToBase64(tt.args.r)); strings.Compare(string(got), tt.want) != 0 {
				t.Errorf("ToBase64() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestFromBase64(t *testing.T) {
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
			args: args{r: strings.NewReader("VGhpcyBpcyBvbmx5IGEgdGVzdCBvZiB0aGUgVG9CYXNlNjQgZmlsdGVyLg==")},
			want: "This is only a test of the ToBase64 filter.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := io.ReadAll(FromBase64(tt.args.r)); strings.Compare(string(got), tt.want) != 0 {
				t.Errorf("FromBase64() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
