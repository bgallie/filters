// Copyright 2020 Billy G. Allie.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package base64 defines filters to encode/decode data base64 encoding.
// These filters can be connected to other filters via io.Pipes.
package base64

import (
	"encoding/base64"
	"io"

	"github.com/bgallie/filters"
)

// ToBase64 reads data from r, encodes it using a base64 encoder.
// The base64 encoded data can be read using the returned PipeReader.
func ToBase64(r io.Reader) *io.PipeReader {
	rRdr, rWrtr := io.Pipe()
	base64W := base64.NewEncoder(base64.StdEncoding, rWrtr)

	go func() {
		defer rWrtr.Close()
		defer base64W.Close()
		_, err := io.Copy(base64W, r)
		filters.CheckFatal(err)
		return
	}()

	return rRdr
}

// FromBase64 reads ascii85 encoded data from r, decodes it using the base64
// decoder.  The decoded data can be read using the returned PipeReader.
func FromBase64(r io.Reader) *io.PipeReader {
	rRdr, rWrtr := io.Pipe()
	base64R := base64.NewDecoder(base64.StdEncoding, r)

	go func() {
		defer rWrtr.Close()
		_, err := io.Copy(rWrtr, base64R)
		filters.CheckFatal(err)
		return
	}()

	return rRdr
}
