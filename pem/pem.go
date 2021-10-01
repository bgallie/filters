// Copyright 2020 Billy G. Allie.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pem defines filters to encode/decode data to/from a stream of
// binary data.  These filters can be connected to other filters via io.Pipes.
package pem

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/bgallie/filters/base64"
	"github.com/bgallie/filters/lines"
	"github.com/friendsofgo/errors"
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
		_, err := fmt.Fprintf(rWrtr, "-----BEGIN %s-----\n", blk.Type)
		if err != nil {
			rWrtr.CloseWithError(errors.Wrap(err, "failure printing PEM BEGIN line to a pipe writer"))
		}
		for k, v := range blk.Headers {
			_, err = fmt.Fprintf(rWrtr, "%s: %s\n", k, v)
			if err != nil {
				rWrtr.CloseWithError(errors.Wrap(err, "failure printing PEM Header line to a pipe writer"))
			}
		}
		lines.LineSize = 76 // output 76 characters per line
		_, err = io.Copy(rWrtr, lines.SplitToLines(base64.ToBase64(r)))
		if err != nil {
			rWrtr.CloseWithError(errors.Wrap(err, "failure copying (io.Copy) base64 encoded data to a pipe writer"))
		}
		_, err = fmt.Fprintf(rWrtr, "-----END %s-----\n", blk.Type)
		if err != nil {
			rWrtr.CloseWithError(errors.Wrap(err, "failure printing PEM END line to a pipe writer"))
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
		base64W.CloseWithError(errors.Wrap(err, "Missing PEM message."))
	}
	// Get the type of PEM message
	if bytes.HasPrefix(line, []byte("-----BEGIN ")) {
		i := bytes.Index(line, []byte(" ")) + 1
		j := bytes.Index(line[i:], []byte("-"))
		blk.Type = string(line[i : i+j])
	} else {
		base64W.CloseWithError(errors.Wrap(err, "Incorrectly formed PEM message: no BEGIN line."))
	}
	// Get the header data if any.
	for {
		line, err = bRdr.ReadBytes('\n')
		if err != nil {
			base64W.CloseWithError(errors.Wrap(err, "Incomplete/malformed PEM message"))
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
				base64W.CloseWithError(errors.Wrap(err, "Failed to write to the base64 encoder"))
			}
			line, err = bRdr.ReadBytes('\n')
			if err != nil {
				base64W.CloseWithError(errors.Wrap(err, "Incomplete/malformed PEM message"))
			}
		}
		if err == nil {
			if bytes.HasPrefix(line, []byte("-----END ")) {
				i := bytes.Index(line, []byte(" ")) + 1
				j := bytes.Index(line[i:], []byte("-"))
				if blk.Type != string(line[i:i+j]) {
					base64W.CloseWithError(errors.New("Incorrectly formed PEM message: BEGIN/END type mismatch"))
				}
			} else {
				base64W.CloseWithError(errors.New("Incorrectly formed PEM message: missing END line"))
			}
		} else {
			base64W.CloseWithError(errors.Wrap(err, "Incorrectly formed PEM message."))
		}
	}()

	return base64.FromBase64(base64R), blk
}
