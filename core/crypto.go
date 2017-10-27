package core

import (
	"crypto/sha256"
	"encoding/binary"
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
// Returns a success, hash pair.
func ProofOfWork(blockHash []byte, counter uint32) (bool, []byte) {
	// convert counter integer into array of bytes
	counterBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(counterBytes, counter)

	// concatenate blockHash and counter bytes
	var bs []byte
	for _, b := range blockHash {
		bs = append(bs, b)
	}
	for _, b := range counterBytes {
		bs = append(bs, b)
	}

	// successful Proof of Work if hashed concatenated bytes ends in
	// number of zeroes given by _difficulty_
	difficulty := 1
	hash := Sha256(bs)
	tail := hash[len(hash)-difficulty:]
	for i := 0; i < difficulty; i++ {
		if tail[i] != 0 {
			return false, hash
		}
	}

	return true, hash
}

// func VerifyProof(lastProof Proof, proof Proof) bool {
// 	guess := fmt.Sprintf("%d%d", lastProof, proof)
// 	guessHash := Sha256([]byte(guess))
// 	return guessHash[:4] == "0000"
// }
