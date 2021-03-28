// Package filters - flate: compress/uncompress a data stream.
package filters

import (
	"compress/flate"
	"io"
)

// ToFlate reads data from r and compresses it using flate with the best
// compression method available to it.  The compressed data can be read using
// the returned PipeReader.
func ToFlate(r io.Reader) *io.PipeReader {
	rRdr, rWrtr := io.Pipe()
	flateW, err := flate.NewWriter(rWrtr, flate.BestCompression)
	checkFatal(err)

	go func() {
		defer rWrtr.Close()
		defer flateW.Close()
		_, err := io.Copy(flateW, r)
		checkFatal(err)
		checkFatal(flateW.Flush())
		return
	}()

	return rRdr
}

// FromFlate reads data compressed using flate from r and decompresses it.
// The decompressed data can be read from the returned PipeReader.
func FromFlate(r io.Reader) *io.PipeReader {
	rRdr, rWrtr := io.Pipe()
	flateR := flate.NewReader(r)

	go func() {
		defer flateR.Close()
		defer rWrtr.Close()
		_, err := io.Copy(rWrtr, flateR)
		checkFatal(err)
		return
	}()

	return rRdr
}
