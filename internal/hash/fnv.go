package hash

import (
	"encoding/hex"
	"hash/fnv"
)

// FNV returns the FNV-1a hash of the given text.
func FNV(text string) uint64 {
	fnvHash := fnv.New64a()
	fnvHash.Write([]byte(text))

	return fnvHash.Sum64()
}

// FNVString returns the FNV-1a hash of the given text as a string.
func FNVString(text string) string {
	fnvHash := fnv.New64a()
	fnvHash.Write([]byte(text))

	return hex.EncodeToString(fnvHash.Sum(nil))
}
