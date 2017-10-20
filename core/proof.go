package core

import "fmt"

func ProofOfWork(lastProof Proof) Proof {
	var proof Proof
	proof = 0
	for !VerifyProof(lastProof, proof) {
		proof++
	}
	return proof
}

func VerifyProof(lastProof Proof, proof Proof) bool {
	guess := fmt.Sprintf("%d%d", lastProof, proof)
	guessHash := Sha256([]byte(guess))
	return guessHash[:4] == "0000"
}
