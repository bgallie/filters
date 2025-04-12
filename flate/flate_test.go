// Copyright 2020 Billy G. Allie.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package flate defines filters to compress/uncompress data using flate.
// These filters can be connected to other filters via io.Pipes.
package flate

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

// formatByteSlice will take a byte slice and format a string
// representation of it.  It will consiste of lines of 16 hexidecimal
// characters seperated by a ', '.
func formatByteSlice(prefix string, src []byte) string {
	var output bytes.Buffer
	output.WriteString(prefix)
	i := 0
	for _, v := range src {
		output.WriteString(fmt.Sprintf("%#02x", v))
		i++
		switch {
		case i == len(src):
			// Do nothine.
		case i%16 == 0:
			output.WriteString(",\n" + prefix)
		default:
			output.WriteString(", ")
		}
	}
	return output.String()
}

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
				0x1c, 0xc9, 0xc1, 0x09, 0xc0, 0x20, 0x0c, 0x40, 0xd1, 0x55, 0xfe, 0x04, 0x19, 0xc3, 0x09, 0x5c,
				0x40, 0xda, 0x88, 0x01, 0x31, 0x60, 0x72, 0xd0, 0xed, 0x0b, 0x3d, 0xbf, 0x3a, 0x2c, 0xb0, 0xc0,
				0xd7, 0xbc, 0x34, 0x52, 0x23, 0xf1, 0x4e, 0x0e, 0xa5, 0x7a, 0x99, 0x2d, 0x95, 0x6e, 0x33, 0x75,
				0x0b, 0x14, 0xdf, 0x3f, 0x2c, 0x3d, 0x49, 0xd8, 0xc9, 0x4b, 0xe8, 0xe3, 0xeb, 0x0d, 0x44, 0xe4,
				0x0b, 0x00, 0x00, 0xff, 0xff,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ttfRdr := bufio.NewReader(ToFlate(tt.args.r))
			if got, _ := io.ReadAll(ttfRdr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToFlate() = %v, want %v", formatByteSlice("", got), formatByteSlice("", tt.want))
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
				0x1c, 0xc9, 0xc1, 0x09, 0xc0, 0x20, 0x0c, 0x40, 0xd1, 0x55, 0xfe, 0x04, 0x19, 0xc3, 0x09, 0xba,
				0x80, 0xb4, 0x11, 0x03, 0xd6, 0x80, 0xc9, 0x41, 0xb7, 0x2f, 0xf4, 0xfc, 0xae, 0x6e, 0x81, 0x05,
				0x3e, 0xc7, 0xa1, 0x92, 0x1a, 0x89, 0x37, 0xb2, 0x2b, 0x65, 0xf9, 0x5b, 0x46, 0x4d, 0xa5, 0xd9,
				0x48, 0x5d, 0x02, 0xc5, 0xd7, 0x4f, 0x53, 0x77, 0x12, 0xb6, 0xf3, 0x10, 0x7a, 0xfb, 0x7c, 0x02,
				0x11, 0xf9, 0x02, 0x00, 0x00, 0xff, 0xff,
			})},
			want: "This is only a test of the FromFlate filter.  For the next sixty seconds ...",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tffRdr := FromFlate(tt.args.r)
			got, _ := io.ReadAll(tffRdr)
			if string(got) != tt.want {
				t.Errorf("FromFlate() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
