// ascii85
package filters

import (
	"encoding/ascii85"
	"io"
)

// ToAscii85 reads data from r, encodes it using Ascii85.
// The Ascii85 encoded data can be read using the returned PipeReader.
func ToAscii85(r io.Reader) *io.PipeReader {
	defer un(trace("ToAscii85"))
	rRdr, rWrtr := io.Pipe()
	ascii85W := ascii85.NewEncoder(rWrtr)

	go func() {
		defer un(trace("ToAscii85 -> encoding ascii85"))
		defer rWrtr.Close()
		defer ascii85W.Close()
		_, err := io.Copy(ascii85W, r)
		checkFatal(err)
		return
	}()

	return rRdr
}

// FromAscii85r reads ascii85 encoded data from r, decodes it using the ascii85
// decoder.  The decoded data can be read using the returned PipeReader.
func FromAscii85(r io.Reader) *io.PipeReader {
	defer un(trace("FromAscii85"))
	rRdr, rWrtr := io.Pipe()
	ascii85R := ascii85.NewDecoder(r)

	go func() {
		defer un(trace("FromAscii85 -> decoding ascii85"))
		defer rWrtr.Close()
		_, err := io.Copy(rWrtr, ascii85R)
		checkFatal(err)
		return
	}()

	return rRdr
}
