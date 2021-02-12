// Package filters - ascii85: convert a stream of binary data to/from ASCII86 encoding.
package filters

import (
	"encoding/ascii85"
	"io"
)

// ToASCII85 reads data from r, encodes it using Ascii85.
// The Ascii85 encoded data can be read using the returned PipeReader.
func ToASCII85(r io.Reader) *io.PipeReader {
	defer un(trace("ToASCII85"))
	rRdr, rWrtr := io.Pipe()
	ascii85W := ascii85.NewEncoder(rWrtr)

	go func() {
		defer un(trace("ToASCII85 -> encoding ascii85"))
		defer rWrtr.Close()
		defer ascii85W.Close()
		_, err := io.Copy(ascii85W, r)
		checkFatal(err)
		return
	}()

	return rRdr
}

// FromASCII85 reads ascii85 encoded data from r, decodes it using the ascii85
// decoder.  The decoded data can be read using the returned PipeReader.
func FromASCII85(r io.Reader) *io.PipeReader {
	defer un(trace("FromASCII85"))
	rRdr, rWrtr := io.Pipe()
	ascii85R := ascii85.NewDecoder(r)

	go func() {
		defer un(trace("FromASCII85 -> decoding ascii85"))
		defer rWrtr.Close()
		_, err := io.Copy(rWrtr, ascii85R)
		checkFatal(err)
		return
	}()

	return rRdr
}
