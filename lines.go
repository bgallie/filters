// lines.go
package filters

import (
	"bufio"
	"fmt"
	"io"
	"log"
)

// SplitToLines reads a stream of ASCII characters (usally the output from
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
			log.Printf("Split2Lines -> io.ReadFull: len: %d, err: %v\n", n, err)
			checkFatal(err)

			if err != nil {
				// must be EOF or ErrUnexpectedEOF
				if n != 0 {
					n, err = fmt.Fprintln(rWtr, string(line[:n]))
					log.Printf("Split2Lines -> io.PipeWriter: len: %d, err: %v\n",
						n, err)
					checkFatal(err)
				}

				break
			}

			n, err = fmt.Fprintln(rWtr, string(line[:n]))
			log.Printf("Split2Lines -> io.PipeWriter: len: %d, err: %v\n", n, err)
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
			line, isPrefix, err := bRdr.ReadLine()
			log.Printf("CombineLines -> ReadLine: len: %d, isPrefix: %v, err: %v\n",
				len(line), isPrefix, err)

			if err == nil {
				n, err := rWtr.Write(line)
				log.Printf("CombineLines -> io.PipeWriter: len: %d, err: %v\n",
					n, err)
				checkFatal(err)
			} else {
				checkFatal(err)
				break
			}
		}
	}()

	return rRdr
}
