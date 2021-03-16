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

func ToPem(r io.Reader, blk Block) *io.PipeReader {
	defer un(trace("ToPem"))
	rRdr, rWrtr := io.Pipe()

	go func() {
		defer un(trace("ToPem -> encoding pem"))
		defer rWrtr.Close()
		fmt.Fprintf(rWrtr, "-----BEGIN %s-----\n", blk.Type)
		for k, v := range blk.Headers {
			fmt.Fprintf(rWrtr, "%s:%s\n", k, v)
		}
		fmt.Fprintln(rWrtr, "")
		_, err := io.Copy(rWrtr, SplitToLines(ToBase64(r)))
		checkFatal(err)
		fmt.Fprintf(rWrtr, "-----END %s-----\n", blk.Type)

		return
	}()

	return rRdr
}

func FromPem(r io.Reader) (*io.PipeReader, Block) {
	var blk Block
	rRdr, rWtr := io.Pipe()
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
	for len(line) != 0 {
		line, _, err := bRdr.ReadLine()
		checkFatal(err)
		i := bytes.Index(line, []byte(":"))
		k := string(line[:i])
		v := string(line[i+1:])
		blk.Headers[k] = v
	}

	go func() {
		defer un(trace("FromPem -> decoding base64"))
		defer rWtr.Close()

		_, err := io.Copy(rWtr, FromBase64(bRdr))
		checkFatal(err)
		line, _, err := bRdr.ReadLine()
		checkFatal(err)
		if err != nil { // must be EOF or UnexpectedEOF
			panic("Incorrectly formed PEM message.\n")
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
		return
	}()

	return rRdr, blk
}
