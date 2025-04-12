// Copyright 2020 Billy G. Allie.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package zlib defines filters to compress/uncompress data using zlib.
// These filters can be connected to other filters via io.Pipes.
package zlib

import (
	"compress/zlib"
	"fmt"
	"io"
	"log"
)

// Tozlib reads data from r and compresses it using zlib with the best
// compression method available to it.  The compressed data can be read using
// the returned PipeReader.
func ToZlib(r io.Reader) *io.PipeReader {
	rRdr, rWrtr := io.Pipe()
	zlibW, err := zlib.NewWriterLevel(rWrtr, zlib.BestCompression)
	if err != nil {
		log.Fatalln(fmt.Errorf("error getting a zlib.NewWriterLevel in ToZlib: %w", err))
	}

	go func() {
		defer rWrtr.Close()
		defer zlibW.Close()
		_, err := io.Copy(zlibW, r)
		if err != nil {
			log.Fatalln(fmt.Errorf("error copying (io.Copy) from an io.Reader to a zlib.Writer: %w", err))
		}
	}()

	return rRdr
}

// FromZlib reads data compressed using zlib from r and decompresses it.
// The decompressed data can be read from the returned PipeReader.
func FromZlib(r io.Reader) *io.PipeReader {
	rRdr, rWrtr := io.Pipe()
	zlibR, err := zlib.NewReader(r)
	if err != nil {
		rWrtr.Close()
		zlibR.Close()
		log.Fatalln(fmt.Errorf("error creating a zlib.NewReader in FromZlib: %w", err))
	}

	go func() {
		defer rWrtr.Close()
		defer zlibR.Close()
		_, err = io.Copy(rWrtr, zlibR)
		if err != nil {
			log.Fatalln(fmt.Errorf("error copying (io.Copy) from a zlib.Reader to an io.PipeWriter: %w", err))
		}
	}()

	return rRdr
}
