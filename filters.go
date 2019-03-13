// filters project filters.go
package filters

import (
	"compress/flate"
	"encoding/ascii85"
	"io"
	"log"

	"gitthub.com/bgallie/utilities"
)

var (
	un         = utilities.Un
	trace      = utilities.Trace
	deferClose = utilities.DeferClose
	checkFatal = utilities.CheckFatal
)

func ToAscii85(r io.Reader) *io.PipeReader {
	defer utilities.Un(trace("ToAscii85"))
	rRdr, rWrtr := io.Pipe()
	ascii85W := ascii85.NewEncoder(rWrtr)

	go func() {
		defer un(trace("ToAscii85 -> encoding ascii85"))
		defer rWrtr.Close()
		defer ascii85W.Close()
		wcnt, err := io.Copy(ascii85W, r)
		log.Printf("ToAscii85 -> encoding ascii85 -> io.Copy wrote %d bytes.  err: %v\n", wcnt, err)
		return
	}()

	return rRdr
}

func FromAscii85(r io.Reader) *io.PipeReader {
	defer un(trace("FromAscii85"))
	rRdr, rWrtr := io.Pipe()
	ascii85R := ascii85.NewDecoder(r)

	go func() {
		defer un(trace("FromAscii85 -> decoding ascii85"))
		defer rWrtr.Close()
		wcnt, err := io.Copy(rWrtr, ascii85R)
		log.Printf("FromAscii85 -> decoding ascii85: io.Copy wrote %d bytes.  err: %v\n", wcnt, err)
		return
	}()

	return rRdr

}

func ToFlate(r io.Reader) *io.PipeReader {
	defer un(trace("ToFlate:"))
	rRdr, rWrtr := io.Pipe()
	flateW, err := flate.NewWriter(rWrtr, flate.BestCompression)
	checkFatal(err)

	go func() {
		defer un(trace("ToFlate -> writing flate"))
		defer rWrtr.Close()
		defer flateW.Close()
		wcnt, err := io.Copy(flateW, r)
		log.Printf("ToFlate -> writing flate: io.Copy wrote %d bytes.  err: %v\n", wcnt, err)
		checkFatal(flateW.Flush())
		return
	}()

	return rRdr
}

func FromFlate(r io.Reader) *io.PipeReader {
	defer un(trace("FromFlate:"))
	rRdr, rWrtr := io.Pipe()
	flateR := flate.NewReader(r)

	go func() {
		defer un(trace("FromFlate -> reading flate"))
		defer flateR.Close()
		defer rWrtr.Close()
		wcnt, err := io.Copy(rWrtr, flateR)
		log.Printf("FromFlate -> reading flate: io.Copy wrote %d bytes.  err: %v\n", wcnt, err)
		return
	}()

	return rRdr
}
