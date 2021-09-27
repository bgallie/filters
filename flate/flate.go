// Copyright 2020 Billy G. Allie.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package flate defines filters to compress/uncompress data using flate.
// These filters can be connected to other filters via io.Pipes.
package flate

import (
	"compress/flate"
	"io"

	"github.com/friendsofgo/errors"
)

// ToFlate reads data from r and compresses it using flate with the best
// compression method available to it.  The compressed data can be read using
// the returned PipeReader.
func ToFlate(r io.Reader) *io.PipeReader {
	rRdr, rWrtr := io.Pipe()
	flateW, err := flate.NewWriter(rWrtr, flate.BestCompression)
	if err != nil {
		rWrtr.CloseWithError(errors.Wrap(err, "failure creating flate.NewWriter."))
	}

	go func() {
		defer rWrtr.Close()
		defer flateW.Close()
		_, err := io.Copy(flateW, r)
		if err != nil {
			rWrtr.CloseWithError(errors.Wrap(err, "failure copying (io.Copy) from a reader to a flate writer."))
		}
	}()

	return rRdr
}

// FromFlate reads data compressed using flate from r and decompresses it.
// The decompressed data can be read from the returned PipeReader.
func FromFlate(r io.Reader) *io.PipeReader {
	rRdr, rWrtr := io.Pipe()
	flateR := flate.NewReader(r)

	go func() {
		defer rWrtr.Close()
		defer flateR.Close()
		_, err := io.Copy(rWrtr, flateR)
		if err != nil {
			rWrtr.CloseWithError(errors.Wrap(err, "failure copying (io.Copy) from a flate reader to a pipe writer."))
		}
	}()

	return rRdr
}
