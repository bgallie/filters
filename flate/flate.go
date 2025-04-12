// Copyright 2020 Billy G. Allie.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package flate defines filters to compress/uncompress data using flate.
// These filters can be connected to other filters via io.Pipes.
package flate

import (
	"compress/flate"
	"fmt"
	"io"
	"log"
)

// ToFlate reads data from r and compresses it using flate with the best
// compression method available to it.  The compressed data can be read using
// the returned PipeReader.
func ToFlate(r io.Reader) *io.PipeReader {
	rRdr, rWrtr := io.Pipe()
	flateW, err := flate.NewWriter(rWrtr, flate.BestCompression)
	if err != nil {
		rRdr.Close()
		rWrtr.Close()
		log.Fatalln(fmt.Errorf("error creating flate.NewWriter: %w", err))
	}

	go func() {
		defer rWrtr.Close()
		defer flateW.Close()
		_, err := io.Copy(flateW, r)
		if err != nil {
			log.Fatalln(fmt.Errorf("error copying to the flate.Writer from an io.Reader: %w", err))
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
		if err != nil && err != io.ErrUnexpectedEOF {
			log.Fatalln(fmt.Errorf("error copying (io.Copy) from a flate.Reader to a io.PipeWriter: %w", err))
		}
	}()

	return rRdr
}
