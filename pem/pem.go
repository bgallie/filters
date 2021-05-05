package pem

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/bgallie/filters"
	"github.com/bgallie/filters/base64"
	"github.com/bgallie/filters/lines"
)

// A Block represents a PEM encoded structure.
//
// The encoded form is:
//    -----BEGIN Type-----
//    Headers
//    base64-encoded Bytes
//    -----END Type-----
// where Headers is a possibly empty sequence of Key: Value lines.
type Block struct {
	Type    string            // The type, taken from the preamble (i.e. "RSA PRIVATE KEY").
	Headers map[string]string // Optional headers.
}

// ToPem reads data from r, encodes it using PEM formatted encoding.
// The PEM encoded data can be read using the returned PipeReader.
func ToPem(r io.Reader, blk Block) *io.PipeReader {
	rRdr, rWrtr := io.Pipe()

	go func() {
		defer rWrtr.Close()
		fmt.Fprintf(rWrtr, "-----BEGIN %s-----\n", blk.Type)
		for k, v := range blk.Headers {
			fmt.Fprintf(rWrtr, "%s: %s\n", k, v)
		}
		lines.LineSize = 64
		_, err := io.Copy(rWrtr, lines.SplitToLines(base64.ToBase64(r)))
		filters.CheckFatal(err)
		fmt.Fprintf(rWrtr, "-----END %s-----\n", blk.Type)

		return
	}()

	return rRdr
}

// FromPem reads PEM encoded data from r, decodes it using a base64
// decoder.  The PEM information is returned in the pem.Block structure
// and the decoded data can be read using the returned PipeReader.
func FromPem(r io.Reader) (*io.PipeReader, Block) {
	var blk Block
	blk.Headers = make(map[string]string)
	base64R, base64W := io.Pipe()
	bRdr := bufio.NewReader(r)
	line, _, err := bRdr.ReadLine()
	filters.CheckFatalMsg(err, "Missing PEM message.")
	// Get the type of PEM message
	if bytes.HasPrefix(line, []byte("-----BEGIN ")) {
		i := bytes.Index(line, []byte(" ")) + 1
		j := bytes.Index(line[i:], []byte("-"))
		blk.Type = string(line[i : i+j])
	} else {
		log.Fatalln("Incorrectly formed PEM message: no BEGIN line.")
	}

	go func() {
		defer base64W.Close()
		// Read the base64 data from the PEM message and send it to the base64
		// filter for decoding after processing any header information.
		line, _, err := bRdr.ReadLine()
		filters.CheckFatalMsg(err, "Incomplete/malformed PEM message")
		for {
			i := bytes.Index(line, []byte(": "))
			if i < 0 {
				break
			}
			k := string(line[:i])
			v := string(line[i+2:])
			blk.Headers[k] = v
			line, _, err = bRdr.ReadLine()
			filters.CheckFatalMsg(err, "Incomplete/malformed PEM message")
		}
		for err == nil && !bytes.HasPrefix(line, []byte("-----END ")) {
			_, err = base64W.Write(line)
			filters.CheckFatal(err)
			line, _, err = bRdr.ReadLine()
			filters.CheckFatalMsg(err, "Incomplete/malformed PEM message")
		}
		if err == nil && bytes.HasPrefix(line, []byte("-----END ")) {
			i := bytes.Index(line, []byte(" ")) + 1
			j := bytes.Index(line[i:], []byte("-"))
			if blk.Type != string(line[i:i+j]) {
				log.Fatalln("Incorrectly formed PEM message: BEGIN/END type mismatch")
			}
		} else {
			log.Fatalln("Incorrectly formed PEM message.")
		}
	}()

	return base64.FromBase64(base64R), blk
}
