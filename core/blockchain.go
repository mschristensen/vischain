package core

import (
	"time"
)

type Blockchain []Block

func NewBlockchain() Blockchain {
	var bc Blockchain
	genesis := Block{
		index:        0,
		timestamp:    time.Now().UnixNano(),
		transactions: nil,
		proof:        0,
		PrevHash:     "Genesis",
	}
	bc = append(bc, genesis)
	return bc
}

func (bc *Blockchain) AddBlock(proof Proof, prevHash string, transactions []Transaction) Block {
	block := Block{
		index:        1,
		timestamp:    time.Now().UnixNano(),
		transactions: transactions,
		proof:        proof,
		PrevHash:     prevHash,
	}
	*bc = append(*bc, block)
	return block
}
