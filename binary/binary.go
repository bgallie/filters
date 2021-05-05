// Copyright 2020 Billy G. Allie.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package binary defines filters to encode the data as a stream of '0' and '1'
// characters.  These filters can be connected to other filters via io.Pipes.
package binary

import (
	"fmt"
	"io"

	"github.com/bgallie/filters"
)

// ToBinary reads data from r, encodes it as a stream of '0' and '1' characters.
// The ToBinary encoded data can be read using the returned PipeReader.
func ToBinary(r io.Reader) *io.PipeReader {
	rRdr, rWrtr := io.Pipe()

	go func() {
		defer rWrtr.Close()
		for {
			buf := make([]byte, 1024)
			cnt, err := r.Read(buf)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			filters.CheckFatal(err)
			cnt *= 8
			for i := 0; i < cnt; i++ {
				if filters.GetBit(buf, uint(i)) {
					fmt.Fprint(rWrtr, "1")
				} else {
					fmt.Fprint(rWrtr, "0")
				}
			}
		}
		return
	}()

	return rRdr
}

// FromBinary reads data encoded by ToBinary from r, and decodes it.
// The decoded data can be read using the returned PipeReader.
func FromBinary(r io.Reader) *io.PipeReader {
	rRdr, rWrtr := io.Pipe()
	buf := make([]byte, 1024)

	go func() {
		defer rWrtr.Close()
		n, err := r.Read(buf)
		if err == io.EOF {
			return
		}
		if err != nil && err != io.ErrUnexpectedEOF {
			filters.CheckFatal(err)
		}
		for {
			outb := make([]byte, 128)

			for i := 0; i < n; i++ {
				switch string(buf[i]) {
				case "1":
					outb = filters.SetBit(outb, uint(i))
				case "0":
					outb = filters.ClrBit(outb, uint(i))
				default:
					panic("Invalid input to FromBinary")
				}
			}

			_, err = rWrtr.Write(outb[:(n+7)/8])
			filters.CheckFatal(err)
			n, err = r.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil && err != io.ErrUnexpectedEOF {
				filters.CheckFatal(err)
			}
		}
		return
	}()

	return rRdr
}