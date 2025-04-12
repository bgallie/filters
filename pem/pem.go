// Copyright 2020 Billy G. Allie.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pem defines filters to encode/decode data to/from a PEM encoded stream of
// binary data.  PEM data encoding originated in Privacy Enhanced Mail. The most
// common use of PEM encoding today is in TLS keys and certificates.  These filters
// can be connected to other filters via io.Pipes.
package pem

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/bgallie/filters/base64"
	"github.com/bgallie/filters/lines"
)

// A Block represents a PEM encoded structure.
//
// The encoded form is:
//
//	-----BEGIN Type-----
//	Headers
//	base64-encoded Bytes
//	-----END Type-----
//
// where Headers is a possibly empty sequence of Key: Value lines.
type Block struct {
	Type    string            // The type, taken from the preamble (i.e. "RSA PRIVATE KEY").
	Headers map[string]string // Optional headers.
}

// ToPem reads data from r, encodes it using PEM formatted encoding.
// The 'blk' parameter provides the "Type" and "Headers" in the encoded form.
// The PEM encoded data can be read using the returned PipeReader.
func ToPem(r io.Reader, blk Block) *io.PipeReader {
	rRdr, rWrtr := io.Pipe()

	go func() {
		defer rWrtr.Close()
		_, err := fmt.Fprintf(rWrtr, "-----BEGIN %s-----\n", blk.Type)
		if err != nil {
			rWrtr.CloseWithError(fmt.Errorf("failure printing PEM BEGIN line to an io.PipeWriter: %w", err))
		}
		for k, v := range blk.Headers {
			_, err = fmt.Fprintf(rWrtr, "%s: %s\n", k, v)
			if err != nil {
				rWrtr.CloseWithError(fmt.Errorf("failure printing PEM Header line to an io.PipeWriter: %w", err))
			}
		}
		lines.LineSize = 76 // output 76 characters per line
		_, err = io.Copy(rWrtr, lines.SplitToLines(base64.ToBase64(r)))
		if err != nil {
			rWrtr.CloseWithError(fmt.Errorf("failure copying (io.Copy) base64 encoded data to an io.PipeWriterr: %w", err))
		}
		_, err = fmt.Fprintf(rWrtr, "-----END %s-----\n", blk.Type)
		if err != nil {
			rWrtr.CloseWithError(fmt.Errorf("failure printing PEM END line to an io.PipeWriter: %w", err))
		}
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
	line, err := bRdr.ReadBytes('\n')
	if err != nil {
		base64W.CloseWithError(fmt.Errorf("missing PEM message: %w", err))
	}
	// Get the type of PEM message
	if bytes.HasPrefix(line, []byte("-----BEGIN ")) {
		i := bytes.Index(line, []byte(" ")) + 1
		j := bytes.Index(line[i:], []byte("-"))
		blk.Type = string(line[i : i+j])
	} else {
		base64W.CloseWithError(fmt.Errorf("incorrectly formed PEM message: no BEGIN line"))
	}
	// Get the header data if any.
	for {
		line, err = bRdr.ReadBytes('\n')
		if err != nil {
			base64W.CloseWithError(fmt.Errorf("incomplete/malformed PEM message: %w", err))
		}
		i := bytes.Index(line, []byte(": "))
		if i < 0 {
			break
		}
		k := string(line[:i])
		v := string(line[i+2 : len(line)-1])
		blk.Headers[k] = v
	}
	// Process the base64 data and validate the 'END' line.
	go func() {
		defer base64W.Close()
		// Read the base64 data from the PEM message and send it to the base64
		// filter for decoding after processing any header information.
		for err == nil && !bytes.HasPrefix(line, []byte("-----END ")) {
			_, err = base64W.Write(line)
			if err != nil {
				base64W.CloseWithError(fmt.Errorf("failed to write to a base64.Encoder: %w", err))
			}
			line, err = bRdr.ReadBytes('\n')
			if err != nil {
				base64W.CloseWithError(fmt.Errorf("incomplete/malformed PEM message: %w", err))
			}
		}
		if err == nil {
			if bytes.HasPrefix(line, []byte("-----END ")) {
				i := bytes.Index(line, []byte(" ")) + 1
				j := bytes.Index(line[i:], []byte("-"))
				if blk.Type != string(line[i:i+j]) {
					base64W.CloseWithError(fmt.Errorf("incorrectly formed PEM message: BEGIN/END type mismatch"))
				}
			} else {
				base64W.CloseWithError(fmt.Errorf("incorrectly formed PEM message: missing END line"))
			}
		} else {
			base64W.CloseWithError(fmt.Errorf("incorrectly formed PEM message."))
		}
	}()

	return base64.FromBase64(base64R), blk
}
