package hash

import (
	"crypto/rand"
	"fmt"
)

// UUID generates a random UUID according to RFC 4122.
func UUID() string {
	n := 16
	b := make([]byte, n)
	rand.Read(b)

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// UUIDv4 generates a random UUID version 4 according to RFC 4122.
func UUIDv4() string {
	n := 16
	b := make([]byte, n)
	rand.Read(b)

	var _0x40 byte = 0x40

	var _0x80 byte = 0x80

	var _0x0f byte = 0x0f

	var _0x3f byte = 0x3f

	b[6] = (b[6] & _0x0f) | _0x40 // version 4
	b[8] = (b[8] & _0x3f) | _0x80 // variant 1

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
