// Copyright 2020 Billy G. Allie.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package filters defines functions shared by other packages that
// provide filters that can be connected together via io.Pipes.
package filters

import "log"

func CheckFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CheckFatalMsg(err error, msg string) {
	if err != nil {
		log.Fatalln(msg)
	}
}

// SetBit - set bit in a byte array
func SetBit(ary []byte, bit uint) []byte {
	ary[bit>>3] |= (1 << (bit & 7))
	return ary
}

// ClrBit - clear bit in a byte array
func ClrBit(ary []byte, bit uint) []byte {
	ary[bit>>3] &= ^(1 << (bit & 7))
	return ary
}

// GetBit - return the value of a bit in a byte array
func GetBit(ary []byte, bit uint) bool {
	return (ary[bit>>3]&(1<<(bit&7)) != 0)
}
