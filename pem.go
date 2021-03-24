package filters

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type Block struct {
	Type    string            // The type, taken from the preamble (i.e. "RSA PRIVATE KEY").
	Headers map[string]string // Optional headers.
}

// ToPem reads data from r, encodes it using PEM formatted encoding.
// The PEM encoded data can be read using the returned PipeReader.
func ToPem(r io.Reader, blk Block) *io.PipeReader {
	defer un(trace("ToPem"))
	rRdr, rWrtr := io.Pipe()

	go func() {
		defer un(trace("ToPem -> encoding pem"))
		defer rWrtr.Close()
		fmt.Fprintf(rWrtr, "-----BEGIN %s-----\n", blk.Type)
		for k, v := range blk.Headers {
			fmt.Fprintf(rWrtr, "%s: %s\n", k, v)
		}
		fmt.Fprintln(rWrtr, "")
		LineSize = 64
		_, err := io.Copy(rWrtr, SplitToLines(ToBase64(r)))
		checkFatal(err)
		fmt.Fprintf(rWrtr, "-----END %s-----\n", blk.Type)

		return
	}()

	return rRdr
}

// FromPem reads PEM encoded data from r, decodes it using a base64
// decoder.  The PEM iformation is returned in the pem.Block structure
// and the decoded data can be read using the returned PipeReader.
func FromPem(r io.Reader) (*io.PipeReader, Block) {
	var blk Block
	blk.Headers = make(map[string]string)
	rRdr, rWtr := io.Pipe()
	base64R, base64W := io.Pipe()
	bRdr := bufio.NewReader(r)

	line, _, err := bRdr.ReadLine()
	checkFatal(err)
	// Get the type of PEM message
	if bytes.HasPrefix(line, []byte("-----BEGIN ")) {
		i := bytes.Index(line, []byte(" ")) + 1
		j := bytes.Index(line[i:], []byte("-"))
		blk.Type = string(line[i : i+j])
	} else {
		panic("Incorrectly formed PEM message.\n")
	}

	line, _, err = bRdr.ReadLine()
	checkFatal(err)
	for len(line) != 0 {
		i := bytes.Index(line, []byte(": "))
		k := string(line[:i])
		v := string(line[i+2:])
		blk.Headers[k] = v
		line, _, err = bRdr.ReadLine()
		checkFatal(err)
	}

	go func() {
		defer base64W.Close()

		line, _, err := bRdr.ReadLine()
		checkFatal(err)
		for !bytes.HasPrefix(line, []byte("-----END ")) {
			_, err = base64W.Write(line)
			line, _, err = bRdr.ReadLine()
			checkFatal(err)
		}
		if bytes.HasPrefix(line, []byte("-----END ")) {
			i := bytes.Index(line, []byte(" ")) + 1
			j := bytes.Index(line[i:], []byte("-"))
			if blk.Type != string(line[i:i+j]) {
				panic("Incorrectly formed PEM message: BEGIN/END type mismatch\n")
			}
		} else {
			panic("Incorrectly formed PEM message.\n")
		}
	}()

	go func() {
		defer rWtr.Close()
		_, err = io.Copy(rWtr, FromBase64(base64R))
		checkFatal(err)
	}()

	return rRdr, blk
}
