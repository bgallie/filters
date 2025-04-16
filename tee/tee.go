// Copyright 2025 Billy G. Allie.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ascii85 defines filters to encode/decode data using ASCII85 encoding.
// These filters can be connected to other filters via io.Pipes.
package tee

import (
	"fmt"
	"io"
	"log"
	"os"
)

func Tee(rdr io.Reader, wrtr io.Writer) *io.PipeReader {
	rRdr, rWrtr := io.Pipe()
	go func(r io.Reader, w io.Writer) {
		defer rWrtr.Close()
		tRdr := io.TeeReader(r, w)
		b := make([]byte, 2048)
		cnt, err := io.ReadAtLeast(tRdr, b, 2048)
		for err == nil || err == io.ErrUnexpectedEOF {
			_, err = rWrtr.Write(b[:cnt])
			if err != nil {
				log.Fatalln(fmt.Errorf("error writing %d bytes to an io.WriteCloser: %w", cnt, err))
			}
			cnt, err = io.ReadAtLeast(r, b, 2048)
		}
		if err != io.EOF {
			log.Fatalln(fmt.Errorf("error reading from an io.Reader: %w", err))
		}
	}(rdr, wrtr)
	return rRdr
}

func TeeToFile(rdr io.Reader, filename string) *io.PipeReader {
	fWrtr, err := os.Create(filename)
	if err != nil {
		log.Fatalln(fmt.Errorf("error creating file [%s]: %w\n", filename, err))
	}
	rRdr, rWrtr := io.Pipe()
	go func(r io.Reader, w io.Writer) {
		defer rWrtr.Close()
		defer fWrtr.Close()
		tRdr := io.TeeReader(r, w)
		b := make([]byte, 2048)
		cnt, err := io.ReadAtLeast(tRdr, b, 2048)
		for err == nil || err == io.ErrUnexpectedEOF {
			_, err = rWrtr.Write(b[:cnt])
			if err != nil {
				log.Fatalln(fmt.Errorf("error writing %d bytes to an io.Writer: %w", cnt, err))
			}
			cnt, err = io.ReadAtLeast(r, b, 2048)
		}
		if err != io.EOF {
			log.Fatalln(fmt.Errorf("error reading from an io.Reader: %w", err))
		}
	}(rdr, fWrtr)
	return rRdr
}
