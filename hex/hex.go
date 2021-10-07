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

	"github.com/friendsofgo/errors"
)

// ToHex reads data from r, encodes it using hex encoder.  The encoded
// data can be read using the returned PipeReader.
func ToHex(r io.Reader) *io.PipeReader {
	rRdr, rWrtr := io.Pipe()
	hexW := hex.NewEncoder(rWrtr)

	go func() {
		defer rWrtr.Close()
		_, err := io.Copy(hexW, r)
		if err != nil {
			rWrtr.CloseWithError(errors.Wrap(err, "failure copying (io.Copy) from a reader to a hex encoder."))
		}
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
		if err != nil {
			rWrtr.CloseWithError(errors.Wrap(err, "failure copying (io.Copy) from a hex decoder to a pipe writer"))
		}
	}()

	return rRdr
}
