package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// Function to compute SHA-256 checksum of a file
func ComputeChecksum(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
