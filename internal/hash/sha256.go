package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

// SHA256 returns the SHA256 hash of the given text.
func SHA256(text string) string {
	hash := sha256.New()
	hash.Write([]byte(text))

	return hex.EncodeToString(hash.Sum(nil))
}
