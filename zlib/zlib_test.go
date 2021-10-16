// Copyright 2020 Billy G. Allie.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package zlib defines filters to compress/uncompress data using zlib.
// These filters can be connected to other filters via io.Pipes.
package zlib

import (
	"bufio"
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestToZlib(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "ttz1",
			args: args{r: strings.NewReader("This is only a test of the ToZlib filter.  For the next sixty seconds ...")},
			want: []byte{
				120, 218, 28, 200, 203, 9, 192, 32, 12, 128, 225, 85, 254, 9, 50, 70,
				39, 240, 212, 91, 31, 17, 3, 98, 192, 228, 160, 219, 23, 122, 253, 74,
				179, 192, 2, 31, 125, 115, 145, 26, 137, 87, 178, 41, 197, 207, 110, 55,
				213, 122, 234, 20, 56, 124, 254, 62, 116, 37, 97, 43, 55, 161, 143, 143,
				55, 16, 145, 47, 0, 0, 255, 255, 174, 124, 25, 55},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ttfRdr := bufio.NewReader(ToZlib(tt.args.r))
			if got, _ := io.ReadAll(ttfRdr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToZlib() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromZlib(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "tfz1",
			args: args{r: bytes.NewReader([]byte{
				120, 218, 28, 201, 187, 13, 196, 32, 12, 128, 225, 85, 254, 9, 60, 6,
				19, 92, 117, 93, 30, 70, 88, 34, 88, 194, 46, 96, 251, 72, 169, 191,
				95, 179, 192, 2, 31, 125, 115, 144, 26, 137, 87, 178, 41, 101, 250, 243,
				239, 118, 82, 173, 167, 78, 129, 226, 243, 147, 161, 43, 9, 91, 185, 9,
				189, 124, 220, 129, 136, 188, 1, 0, 0, 255, 255, 231, 166, 26, 8},
			)},
			want: "This is only a test of the FromZlib filter.  For the next sixty seconds ...",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tfzRdr := bufio.NewReader(FromZlib(tt.args.r))
			got, _ := io.ReadAll(tfzRdr)
			if string(got) != tt.want {
				t.Errorf("FromZlib() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
