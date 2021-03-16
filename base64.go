package filters

import (
	"encoding/base64"
	"io"
)

func ToBase64(r io.Reader) *io.PipeReader {
	defer un(trace("ToBase64"))
	rRdr, rWrtr := io.Pipe()
	base64W := base64.NewEncoder(base64.StdEncoding, rWrtr)

	go func() {
		defer un(trace("ToBase64 -> encoding base64"))
		defer rWrtr.Close()
		defer base64W.Close()
		_, err := io.Copy(base64W, r)
		checkFatal(err)
		return
	}()

	return rRdr
}

func FromBase64(r io.Reader) *io.PipeReader {
	defer un(trace("FromBase64"))
	rRdr, rWrtr := io.Pipe()
	base64R := base64.NewDecoder(base64.StdEncoding, r)

	go func() {
		defer un(trace("FromBase64 -> decoding base64"))
		defer rWrtr.Close()
		_, err := io.Copy(rWrtr, base64R)
		checkFatal(err)
		return
	}()

	return rRdr
}
