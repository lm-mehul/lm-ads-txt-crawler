package utils

import "crypto/sha256"

// GenerateHash generates a SHA-256 hash of the given byte slice
func GenerateHash(data []byte) []byte {
	hasher := sha256.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}
