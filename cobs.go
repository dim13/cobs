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

// Encode a null-terminated slice of bytes to a cobs frame
func Encode(p []byte) []byte {
	if len(p) == 0 {
		return nil
	}
	// pad inital message with zero, if missing
	if p[len(p)-1] != 0 {
		p = append(p, 0)
	}
	var buf bytes.Buffer
	for {
		i := bytes.IndexByte(p, 0)
		// no more zeros, we are done
		if i < 0 {
			return buf.Bytes()
		}
		// split oversized chunks
		for i >= 254 {
			buf.WriteByte(255)
			buf.Write(p[:254])
			p = p[254:]
			i -= 254
		}
		// write rest of the chunk
		buf.WriteByte(byte(i + 1))
		buf.Write(p[:i])
		p = p[i+1:]
	}
}

// Decode a cobs frame to a null-terminated slice of bytes
func Decode(p []byte) []byte {
	if len(p) == 0 {
		return nil
	}
	var buf bytes.Buffer
	for {
		// nothing left, we are done
		if len(p) == 0 {
			return buf.Bytes()
		}
		n, body := p[0], p[1:]
		// invalid frame, abort
		if int(n-1) > len(body) || n == 0 {
			return nil
		}
		buf.Write(body[:n-1])
		// full blocks are not followed by zero
		if n < 255 {
			buf.WriteByte(0)
		}
		p = p[n:]
	}
}
