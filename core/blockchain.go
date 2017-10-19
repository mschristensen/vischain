package core

import (
	"time"
)

type Block struct {
	index        int64
	timestamp    int64
	transactions []Transaction
	proof        uint32
	PrevHash     string
}

type Blockchain []Block

func (bc *Blockchain) AddBlock(proof uint32, prevHash string) Block {
	block := Block{
		index:        1,
		timestamp:    time.Now().UnixNano(),
		transactions: nil,
		proof:        proof,
		PrevHash:     prevHash,
	}
	*bc = append(*bc, block)
	return block
}

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
