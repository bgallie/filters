// Copyright 2020 Billy G. Allie.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package flate defines filters to compress/uncompress data using flate.
// These filters can be connected to other filters via io.Pipes.
package flate

import (
	"bufio"
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestToFlate(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "ttf1",
			args: args{r: strings.NewReader("This is only a test of the ToFlate filter.  For the next sixty seconds ...")},
			want: []byte{
				28, 201, 193, 9, 192, 32, 12, 64, 209, 85, 254, 4, 25, 195, 9, 92,
				64, 218, 136, 1, 49, 96, 114, 208, 237, 11, 61, 191, 58, 44, 176, 192,
				215, 188, 52, 82, 35, 241, 78, 14, 165, 122, 153, 45, 149, 110, 51, 117,
				11, 20, 223, 63, 44, 61, 73, 216, 201, 75, 232, 227, 235, 13, 68, 228,
				11, 0, 0, 255, 255,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ttfRdr := bufio.NewReader(ToFlate(tt.args.r))
			if got, _ := io.ReadAll(ttfRdr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToFlate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromFlate(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "tff1",
			args: args{r: bytes.NewReader([]byte{
				28, 201, 193, 9, 192, 32, 12, 64, 209, 85, 254, 4, 25, 195, 9, 186,
				128, 180, 17, 3, 214, 128, 201, 65, 183, 47, 244, 252, 174, 110, 129, 5,
				62, 199, 161, 146, 26, 137, 55, 178, 43, 101, 249, 91, 70, 77, 165, 217,
				72, 93, 2, 197, 215, 79, 83, 119, 18, 182, 243, 16, 122, 251, 124, 2,
				17, 249, 2, 0, 0, 255, 255,
			})},
			want: "This is only a test of the FromFlate filter.  For the next sixty seconds ...",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tffRdr := bufio.NewReader(FromFlate(tt.args.r))
			got, _ := io.ReadAll(tffRdr)
			if string(got) != tt.want {
				t.Errorf("FromFlate() = %v, want %v", got, tt.want)
			}
		})
	}
}
