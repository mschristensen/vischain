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

// ProofOfWork implements a HashCash-based PoW algorithm.
// Work is completed successfully if the hash of the block and a
// given counter ends with a sufficient number of trailing zeroes.
func ProofOfWork(blockHash []byte, counter []byte) bool {
	bs := ConcatBytes(blockHash, counter)
	difficulty := 1
	hash := Sha256(bs)
	tail := hash[len(hash)-difficulty:]
	for i := 0; i < difficulty; i++ {
		if tail[i] != 0 {
			return false
		}
	}

	return true
}

func ConcatBytes(a []byte, b []byte) []byte {
	var result []byte
	for _, i := range a {
		result = append(result, i)
	}
	for _, i := range b {
		result = append(result, i)
	}
	return result
}

func CompareHashes(a Hash, b Hash) bool {
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
