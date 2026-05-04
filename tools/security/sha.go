package security

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

// HashSHA256 is a basic helper to sum and encode
// a string hashed in SHA256.
func HashSHA256(v []byte) string {
	sum := sha256.Sum256(v)
	return hex.EncodeToString(sum[:])
}

// HashSHA512 is a basic helper to sum and encode
// a string hashed in SHA512.
func HashSHA512(v []byte) string {
	sum := sha512.Sum512(v)
	return hex.EncodeToString(sum[:])
}
