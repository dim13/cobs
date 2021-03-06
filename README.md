# Consistent Overhead Byte Stuffing

COBS encoder breaks a packet into one or more sequences of non-zero
bytes.  The encoding routine searches through the first 254 bytes of
the packet looking for the first occurrence of a zero byte.  If no
zero is found, then a code of 0xFF is output, followed by the 254
non-zero bytes. If a zero is found, then the number of bytes
examined, n, is output as the code byte, followed by the actual
values of the (n-1) non-zero bytes up to (but not including) the zero
byte. This process is repeated until all the bytes of the packet have
been encoded.

## Links

* http://www.stuartcheshire.org/papers/COBSforToN.pdf
* https://pkg.go.dev/github.com/dim13/cobs
