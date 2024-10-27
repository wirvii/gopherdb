package hash

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5 returns the MD5 hash of the given text.
func MD5(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))

	return hex.EncodeToString(hash.Sum(nil))
}
