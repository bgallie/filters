package filters

import "log"

func checkFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkFatalMsg(err error, msg string) {
	if err != nil {
		log.Fatalln(msg)
	}
}

// SetBit - set bit in a byte array
func setBit(ary []byte, bit uint) []byte {
	ary[bit>>3] |= (1 << (bit & 7))
	return ary
}

// ClrBit - clear bit in a byte array
func clrBit(ary []byte, bit uint) []byte {
	ary[bit>>3] &= ^(1 << (bit & 7))
	return ary
}

// GetBit - return the value of a bit in a byte array
func getBit(ary []byte, bit uint) bool {
	return (ary[bit>>3]&(1<<(bit&7)) != 0)
}
