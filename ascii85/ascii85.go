// Copyright 2020 Billy G. Allie.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ascii85 defines filters to encode/decode data using ASCII85 encoding.
// These filters can be connected to other filters via io.Pipes.
package ascii85

import (
	"encoding/ascii85"
	"io"

	"github.com/bgallie/filters"
)

// ToASCII85 reads data from r, encodes it using Ascii85.
// The Ascii85 encoded data can be read using the returned PipeReader.
func ToASCII85(r io.Reader) *io.PipeReader {
	rRdr, rWrtr := io.Pipe()
	ascii85W := ascii85.NewEncoder(rWrtr)

	go func() {
		defer rWrtr.Close()
		defer ascii85W.Close()
		_, err := io.Copy(ascii85W, r)
		filters.CheckFatal(err)
		return
	}()

	return rRdr
}

// FromASCII85 reads ascii85 encoded data from r, decodes it using the ascii85
// decoder.  The decoded data can be read using the returned PipeReader.
func FromASCII85(r io.Reader) *io.PipeReader {
	rRdr, rWrtr := io.Pipe()
	ascii85R := ascii85.NewDecoder(r)

	go func() {
		defer rWrtr.Close()
		_, err := io.Copy(rWrtr, ascii85R)
		filters.CheckFatal(err)
		return
	}()

	return rRdr
}
