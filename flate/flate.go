// Copyright 2020 Billy G. Allie.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package flate defines filters to compress/uncompress data using flate.
// These filters can be connected to other filters via io.Pipes.
package flate

import (
	"compress/flate"
	"io"

	"github.com/bgallie/filters"
)

// ToFlate reads data from r and compresses it using flate with the best
// compression method available to it.  The compressed data can be read using
// the returned PipeReader.
func ToFlate(r io.Reader) *io.PipeReader {
	rRdr, rWrtr := io.Pipe()
	flateW, err := flate.NewWriter(rWrtr, flate.BestCompression)
	filters.CheckFatal(err)

	go func() {
		defer rWrtr.Close()
		defer flateW.Close()
		_, err := io.Copy(flateW, r)
		filters.CheckFatal(err)
		filters.CheckFatal(flateW.Flush())
		return
	}()

	return rRdr
}

// FromFlate reads data compressed using flate from r and decompresses it.
// The decompressed data can be read from the returned PipeReader.
func FromFlate(r io.Reader) *io.PipeReader {
	rRdr, rWrtr := io.Pipe()
	flateR := flate.NewReader(r)

	go func() {
		defer flateR.Close()
		defer rWrtr.Close()
		_, err := io.Copy(rWrtr, flateR)
		filters.CheckFatal(err)
		return
	}()

	return rRdr
}
