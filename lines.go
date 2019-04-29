// lines.go
package filters

import (
	"bufio"
	"fmt"
	"io"
)

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
