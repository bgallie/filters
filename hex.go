// hex
package filters

import (
	"encoding/hex"
	"io"
)

// ToHex reads data from r, encodes it using hexadecimal.
// The hexadecimal encoded data can be read using the returned PipeReader.
func ToHex(r io.Reader) *io.PipeReader {
	defer un(trace("ToHex"))
	rRdr, rWrtr := io.Pipe()
	hexW := hex.NewEncoder(rWrtr)

	go func() {
		defer un(trace("ToHex -> encoding hexadecimal"))
		defer rWrtr.Close()
		_, err := io.Copy(hexW, r)
		checkFatal(err)
		return
	}()

	return rRdr
}

// FromHex reads hexadecimal encoded data from r, decodes it using the hex
// decoder.  The decoded data can be read using the returned PipeReader.
func FromHex(r io.Reader) *io.PipeReader {
	defer un(trace("FromAscii85"))
	rRdr, rWrtr := io.Pipe()
	hexR := hex.NewDecoder(r)

	go func() {
		defer un(trace("FromHex -> decoding hexadecimal"))
		defer rWrtr.Close()
		_, err := io.Copy(rWrtr, hexR)
		checkFatal(err)
		return
	}()

	return rRdr
}
