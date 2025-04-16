// Copyright 2025 Billy G. Allie.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ascii85 defines filters to encode/decode data using ASCII85 encoding.
// These filters can be connected to other filters via io.Pipes.
package tee

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

type myBuffer struct {
	bufr bytes.Buffer
}

func (b *myBuffer) Write(p []byte) (n int, err error) {
	return b.bufr.Write(p)
}

func (b *myBuffer) CLose() error {
	return nil
}

func TestTee(t *testing.T) {
	var bufr bytes.Buffer
	b.Grow(1024)
	type args struct {
		rdr  io.Reader
		wrtr io.Writer
	}
	tests := []struct {
		name string
		args args
		want *io.PipeReader
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Tee(tt.args.rdr, tt.args.wrtr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tee() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTeeToFile(t *testing.T) {
	type args struct {
		rdr      io.Reader
		filename string
	}
	tests := []struct {
		name string
		args args
		want *io.PipeReader
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TeeToFile(tt.args.rdr, tt.args.filename); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TeeToFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
