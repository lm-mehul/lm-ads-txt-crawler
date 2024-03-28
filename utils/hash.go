package utils

import (
	"crypto/sha256"
	"encoding/base64"
)

// GenerateHash generates a SHA-256 hash of the given byte slice
func GenerateHash(data []byte) string {
	hasher := sha256.New()
	hasher.Write(data)
	hashInBytes := hasher.Sum(nil)
	return base64.StdEncoding.EncodeToString(hashInBytes)
}
