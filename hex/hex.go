// Copyright 2020 Billy G. Allie.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package hex defines filters to encode/decode data to/from a stream of
// hexadecimal characters.  These filters can be connected to other filters
// via io.Pipes.
package hex

import (
	"encoding/hex"
	"io"

	"github.com/bgallie/filters"
)

// ToHex reads data from r, encodes it using hexadecimal.
func ToHex(r io.Reader) *io.PipeReader {
	rRdr, rWrtr := io.Pipe()
	hexW := hex.NewEncoder(rWrtr)

	go func() {
		defer rWrtr.Close()
		_, err := io.Copy(hexW, r)
		filters.CheckFatal(err)
		return
	}()

	return rRdr
}

// FromHex reads hexadecimal encoded data from r, decodes it using the hex
// decoder.  The decoded data can be read using the returned PipeReader.
func FromHex(r io.Reader) *io.PipeReader {
	rRdr, rWrtr := io.Pipe()
	hexR := hex.NewDecoder(r)

	go func() {
		defer rWrtr.Close()
		_, err := io.Copy(rWrtr, hexR)
		filters.CheckFatal(err)
		return
	}()

	return rRdr
}
