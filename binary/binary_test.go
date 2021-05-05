// Copyright 2020 Billy G. Allie.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package binary defines filters to encode the data as a stream of '0' and '1'
// characters.  These filters can be connected to other filters via io.Pipes.
package binary

import (
	"io"
	"strings"
	"testing"
)

func TestToBinary(t *testing.T) {
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
			args: args{r: strings.NewReader("This is only a test of the ToBinary filter.")},
			want: "0010101000010110100101101100111000000100100101101100111000000100" +
				"1111011001110110001101101001111000000100100001100000010000101110" +
				"1010011011001110001011100000010011110110011001100000010000101110" +
				"0001011010100110000001000010101011110110010000101001011001110110" +
				"1000011001001110100111100000010001100110100101100011011000101110" +
				"101001100100111001110100",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := io.ReadAll(ToBinary(tt.args.r)); strings.Compare(string(got), tt.want) != 0 {
				t.Errorf("ToBinary() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestFromBinary(t *testing.T) {
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
			args: args{r: strings.NewReader("0010101000010110100101101100111000000100100101101100111000000100" +
				"1111011001110110001101101001111000000100100001100000010000101110" +
				"1010011011001110001011100000010011110110011001100000010000101110" +
				"0001011010100110000001000110001001001110111101101011011001000010" +
				"1001011001110110100001100100111010011110000001000110011010010110" +
				"0011011000101110101001100100111001110100")},
			want: "This is only a test of the FromBinary filter.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := io.ReadAll(FromBinary(tt.args.r)); strings.Compare(string(got), tt.want) != 0 {
				t.Errorf("FromBinary() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
