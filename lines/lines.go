// Copyright 2020 Billy G. Allie.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package lines defines filters to split a stream of ASCII characters (usually
// the output ascii86) into lines of text and to combine lines of text into
// a stream of ASCII charaters.  These filters can be connected to other filters
// via io.Pipes.
package lines

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

// LineSize defines the number of character that are put into a line by the
// SplitToLines function.  The default is 72 characters per line.  Change
// it prior to calling SplitToLines to change the number of characters
// per line.
var LineSize int = 72

// SplitToLines reads a stream of ASCII characters (usually the output from
// ascii85) from r and splits it into lines of 'LineSize' characters.  The
// lines can be read from the returned PipeReader.
func SplitToLines(r io.Reader) *io.PipeReader {
	rRdr, rWtr := io.Pipe()
	line := make([]byte, LineSize)

	go func() {
		defer rWtr.Close()

		for {
			n, err := io.ReadFull(r, line)
			if err != nil && !(errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF)) {
				rWtr.CloseWithError(fmt.Errorf("error reading a line of text from an io.Reader: %w", err))
			}

			if err != nil {
				// must be EOF or ErrUnexpectedEOF
				if n != 0 {
					_, err = fmt.Fprintln(rWtr, string(line[:n]))
					if err != nil {
						rWtr.CloseWithError(fmt.Errorf("error writing text to an io.PipeWriter: %w", err))
					}
				}

				break
			}

			_, err = fmt.Fprintln(rWtr, string(line[:n]))
			if err != nil {
				rWtr.CloseWithError(fmt.Errorf("error writing text to an io.PipeWriter: %w", err))
			}
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
				if err != nil {
					rWtr.CloseWithError(fmt.Errorf("error writing text to an io.PipeWriter: %w", err))
				}
			} else {
				if !errors.Is(err, io.EOF) {
					rWtr.CloseWithError(fmt.Errorf("error reading a line of text from a buffered io.Reader: %w", err))
				}
				break
			}
		}
	}()

	return rRdr
}
