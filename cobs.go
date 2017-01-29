// Package cobs implements Consistent Overhead Byte Stuffing algorithm
//
// COBS encoder breaks a packet into one or more sequences of non-zero
// bytes.  The encoding routine searches through the first 254 bytes of
// the packet looking for the first occurrence of a zero byte.  If no
// zero is found, then a code of 0xFF is output, followed by the 254
// non-zero bytes. If a zero is found, then the number of bytes
// examined, n, is output as the code byte, followed by the actual
// values of the (n-1) non-zero bytes up to (but not including) the zero
// byte. This process is repeated until all the bytes of the packet have
// been encoded.
//
// See also: http://www.stuartcheshire.org/papers/COBSforToN.pdf
package cobs

import "bytes"

// EncodedSize calculates size of encoded message
func EncodedSize(n int) int {
	return n + n/254 + 1
}

// Encode a slice of bytes to a null-terminated frame
func Encode(p []byte) []byte {
	buf := new(bytes.Buffer)
	writeBlock := func(p []byte) {
		buf.WriteByte(byte(len(p) + 1))
		buf.Write(p)
	}
	for _, ch := range bytes.Split(p, []byte{0}) {
		for len(ch) > 0xfe {
			writeBlock(ch[:0xfe])
			ch = ch[0xfe:]
		}
		writeBlock(ch)
	}
	buf.WriteByte(0)
	return buf.Bytes()
}

// Decode a null-terminated frame to a slice of bytes
func Decode(b []byte) []byte {
	if len(b) == 0 {
		return nil
	}
	buf := new(bytes.Buffer)
	for n := b[0]; n > 0; n = b[0] {
		if int(n) >= len(b) {
			return nil
		}
		buf.Write(b[1:n])
		b = b[n:]
		if n < 0xff && b[0] > 0 {
			buf.WriteByte(0)
		}
	}
	return buf.Bytes()
}
