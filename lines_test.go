// Package filters - lines: split a stream of characters into 72 character lines.
// 							join lines of charaters into a stream of characters.
package filters

import (
	"io"
	"strings"
	"testing"
)

func TestSplitToLines(t *testing.T) {
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
			args: args{r: strings.NewReader("This is only a test of the SplitToLines filter.")},
			want: "This is on\nly a test \nof the Spl\nitToLines \nfilter.\n",
		},
	}
	LineSize = 10
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := io.ReadAll(SplitToLines(tt.args.r)); strings.Compare(string(got), tt.want) != 0 {
				t.Errorf("SplitToLines() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestCombineLines(t *testing.T) {
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
			args: args{r: strings.NewReader("This is on\nly a test \nof the CombineLines \nfilter.\n")},
			want: "This is only a test of the CombineLines filter.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := io.ReadAll(CombineLines(tt.args.r)); strings.Compare(string(got), tt.want) != 0 {
				t.Errorf("CombineLines() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
