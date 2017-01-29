# Consistent Overhead Byte Stuffing

[![Build Status](https://travis-ci.org/dim13/cobs.svg?branch=master)](https://travis-ci.org/dim13/cobs)
[![GoDoc](https://godoc.org/github.com/dim13/cobs?status.svg)](https://godoc.org/github.com/dim13/cobs)

COBS encoder breaks a packet into one or more sequences of non-zero
bytes.  The encoding routine searches through the first 254 bytes of
the packet looking for the first occurrence of a zero byte.  If no
zero is found, then a code of 0xFF is output, followed by the 254
non-zero bytes. If a zero is found, then the number of bytes
examined, n, is output as the code byte, followed by the actual
values of the (n-1) non-zero bytes up to (but not including) the zero
byte. This process is repeated until all the bytes of the packet have
been encoded.

See also: http://www.stuartcheshire.org/papers/COBSforToN.pdf
