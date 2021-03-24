// Package filters - binary: output a stream of bytes as a sequect of '0' and '1' characters.
package filters

import (
	"fmt"
	"io"

	"github.com/bgallie/tnt2/cryptors/bitops"
)

// ToBinary reads data from r, encodes it as a stream of '0' and '1' characters.
// The ToBinary encoded data can be read using the returned PipeReader.
func ToBinary(r io.Reader) *io.PipeReader {
	defer un(trace("ToBinary"))
	rRdr, rWrtr := io.Pipe()

	go func() {
		defer un(trace("ToBinary -> encoding binary characters"))
		defer rWrtr.Close()
		for {
			buf := make([]byte, 1024)
			cnt, err := r.Read(buf)
			checkFatal(err)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			cnt *= 8
			for i := 0; i < cnt; i++ {
				if bitops.GetBit(buf, uint(i)) {
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
// TO DO:  implement FromBinary
// func FromBinary(r io.Reader) *io.PipeReader {
// 	defer un(trace("FromBinary"))
// 	rRdr, rWrtr := io.Pipe()
// 	// inp := make([]byte, 1024)

// 	go func() {
// 		defer un(trace("FromBinary -> decoding binary"))
// 		defer rWrtr.Close()
// 		return
// 	}()

// 	return rRdr
// }
