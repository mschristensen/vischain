package core

import (
	"crypto/sha256"
)

// Sha256 computes the she256 hash of a given byte array
func Sha256(bytes []byte) []byte {
	h := sha256.New()
	h.Write(bytes)
	return h.Sum(nil)
}

func CompareHashes(a Hash, b Hash) bool {
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
