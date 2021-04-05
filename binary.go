// Package filters - binary: output a stream of bytes as a sequect of '0' and '1' characters.
package filters

import (
	"fmt"
	"io"
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
			checkFatal(err)
			cnt *= 8
			for i := 0; i < cnt; i++ {
				if getBit(buf, uint(i)) {
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
			checkFatal(err)
		}
		for {
			outb := make([]byte, 128)

			for i := 0; i < n; i++ {
				switch string(buf[i]) {
				case "1":
					outb = setBit(outb, uint(i))
				case "0":
					outb = clrBit(outb, uint(i))
				default:
					panic("Invalid input to FromBinary")
				}
			}

			_, err = rWrtr.Write(outb[:(n+7)/8])
			checkFatal(err)
			n, err = r.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil && err != io.ErrUnexpectedEOF {
				checkFatal(err)
			}
		}
		return
	}()

	return rRdr
}
