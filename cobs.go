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

// EncodedSize calculates size of encoded message
func EncodedSize(n int) int {
	return n + n/256
}

// Terminate a byte slice with null
func Terminate(p []byte) []byte {
	if len(p) > 0 && p[len(p)-1] == 0 {
		return p
	}
	return append(p, 0)
}

// TrimNull removes null-termination
func TrimNull(p []byte) []byte {
	if len(p) > 0 && p[len(p)-1] == 0 {
		return p[:len(p)-1]
	}
	return p
}

// Encode a null-terminated slice of bytes to a cobs frame
func Encode(p []byte) []byte {
	// preallocate output buffer with estimated size
	b := make([]byte, 0, EncodedSize(len(p)))
	var x [0xff]byte
	x[0] = 1
	for _, v := range p {
		if v == 0 {
			b = append(b, x[:x[0]]...)
			x[0] = 1
		} else {
			x[x[0]] = v
			x[0]++
			if x[0] == 0xff {
				b = append(b, x[:x[0]]...)
				x[0] = 1
			}
		}
	}
	if x[0] > 1 {
		b = append(b, x[:x[0]]...)
	}
	return b
}

// Decode a cobs frame to a null-terminated slice of bytes
func Decode(p []byte) []byte {
	// preallocate output buffer, data will be smaller than input
	b := make([]byte, 0, len(p))
	for len(p) > 0 {
		n, data := p[0], p[1:]
		// invalid frame, abort
		if int(n-1) > len(data) || n == 0 {
			return b
		}
		b = append(b, data[:n-1]...)
		// full blocks are not followed by zero
		if n < 255 {
			b = append(b, 0)
		}
		p = p[n:]
	}
	return b
}
