// Package filters - lines: split a stream of characters into 72 character lines.
// 							join lines of charaters into a stream of characters.
package filters

import (
	"bufio"
	"fmt"
	"io"
)

// SplitToLines reads a stream of ASCII characters (usually the output from
// ascii85) from r and splits it into lines of 72 characters.  The lines can be
// read from the returned PipeReader.
func SplitToLines(r io.Reader) *io.PipeReader {
	defer un(trace("SplitLines:"))
	rRdr, rWtr := io.Pipe()
	var line [72]byte

	go func() {
		defer un(trace("Split2Lines"))
		defer rWtr.Close()

		for {
			n, err := io.ReadFull(r, line[:])
			checkFatal(err)

			if err != nil {
				// must be EOF or ErrUnexpectedEOF
				if n != 0 {
					n, err = fmt.Fprintln(rWtr, string(line[:n]))
					checkFatal(err)
				}

				break
			}

			n, err = fmt.Fprintln(rWtr, string(line[:n]))
			checkFatal(err)
		}
	}()

	return rRdr
}

// CombineLines reads lines of 72 characters (usually the output from
// SplitToLines) and combines them into a stream of characters (minus
// the new line characters).  The stream of characters can be read from
// the returned PipeReader.
func CombineLines(r io.Reader) *io.PipeReader {
	defer un(trace("CombineLines:"))
	rRdr, rWtr := io.Pipe()

	go func() {
		defer un(trace("CombineLines -> PipeWriter:"))
		defer rWtr.Close()
		bRdr := bufio.NewReader(r)

		for {
			line, _, err := bRdr.ReadLine()

			if err == nil {
				_, err := rWtr.Write(line)
				checkFatal(err)
			} else {
				checkFatal(err)
				break
			}
		}
	}()

	return rRdr
}
