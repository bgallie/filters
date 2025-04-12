// Copyright 2020 Billy G. Allie.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package binary defines filters to encode the data as a stream of '0' and '1'
// characters.  These filters can be connected to other filters via io.Pipes.
package binary

import (
	"errors"
	"fmt"
	"io"
)

// SetBit - set bit in a byte array
func SetBit(ary []byte, bit uint) []byte {
	ary[bit>>3] |= (1 << (bit & 7))
	return ary
}

// ClrBit - clear bit in a byte array
func ClrBit(ary []byte, bit uint) []byte {
	ary[bit>>3] &= ^(1 << (bit & 7))
	return ary
}

// GetBit - return the value of a bit in a byte array
func GetBit(ary []byte, bit uint) bool {
	return (ary[bit>>3]&(1<<(bit&7)) != 0)
}

// ToBinary reads data from r, encodes it as a stream of '0' and '1' characters.
// The ToBinary encoded data can be read using the returned PipeReader.
func ToBinary(r io.Reader) *io.PipeReader {
	rRdr, rWrtr := io.Pipe()

	go func() {
		defer rWrtr.Close()
		for {
			buf := make([]byte, 1024)
			cnt, err := r.Read(buf)
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				break
			} else if err != nil {
				rWrtr.CloseWithError(fmt.Errorf("error reading from an io.Reader: %w", err))
			}
			cnt *= 8
			for i := 0; i < cnt; i++ {
				if GetBit(buf, uint(i)) {
					_, err = fmt.Fprint(rWrtr, "1")
				} else {
					_, err = fmt.Fprint(rWrtr, "0")
				}
				if err != nil {
					rWrtr.CloseWithError(fmt.Errorf("error printing to an io.PipeWriter: %w", err))
				}
			}
		}
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
		if errors.Is(err, io.EOF) {
			return
		} else if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
			rWrtr.CloseWithError(fmt.Errorf("error reading from an io.Reader: %w", err))
		}
		for {
			outb := make([]byte, 128)

			for i := 0; i < n; i++ {
				switch string(buf[i]) {
				case "1":
					outb = SetBit(outb, uint(i))
				case "0":
					outb = ClrBit(outb, uint(i))
				default:
					rWrtr.CloseWithError(fmt.Errorf("invalid input to FromBinary"))
				}
			}

			_, err = rWrtr.Write(outb[:(n+7)/8])
			if err != nil {
				rWrtr.CloseWithError(fmt.Errorf("error writing to an io.PipeWriter: %w", err))
			}
			n, err = r.Read(buf)
			if err == io.EOF {
				break
			} else if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
				rWrtr.CloseWithError(fmt.Errorf("error reading from an io.Reader: %w", err))
			}
		}
	}()

	return rRdr
}
