// Copyright 2020 Billy G. Allie.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package hex defines filters to split a stream of ADCII characters (usually
// the output ascii86) into lines of text and to combine lines of text into
// a stream of ASCII charaters.  These filters can be connected to other filters
// via io.Pipes.
package lines

import (
	"bufio"
	"fmt"
	"io"

	"github.com/bgallie/filters"
)

var LineSize int = 72

// SplitToLines reads a stream of ASCII characters (usually the output from
// ascii85) from r and splits it into lines of 'LineSize' characters.  The
// lines can be read from the returned PipeReader.
func SplitToLines(r io.Reader) *io.PipeReader {
	rRdr, rWtr := io.Pipe()
	line := make([]byte, LineSize, LineSize)

	go func() {
		defer rWtr.Close()

		for {
			n, err := io.ReadFull(r, line)
			if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				filters.CheckFatal(err)
			}

			if err != nil {
				// must be EOF or ErrUnexpectedEOF
				if n != 0 {
					n, err = fmt.Fprintln(rWtr, string(line[:n]))
					filters.CheckFatal(err)
				}

				break
			}

			n, err = fmt.Fprintln(rWtr, string(line[:n]))
			filters.CheckFatal(err)
		}
	}()

	return rRdr
}

// CombineLines reads lines of characters (usually the output from
// SplitToLines) and combines them into a stream of characters (minus
// the new line characters).  The stream of characters can be read from
// the returned PipeReader.
func CombineLines(r io.Reader) *io.PipeReader {
	rRdr, rWtr := io.Pipe()

	go func() {
		defer rWtr.Close()
		bRdr := bufio.NewReader(r)

		for {
			line, _, err := bRdr.ReadLine()

			if err == nil {
				_, err := rWtr.Write(line)
				filters.CheckFatal(err)
			} else {
				if err != io.EOF {
					filters.CheckFatal(err)
				}
				break
			}
		}
	}()

	return rRdr
}
