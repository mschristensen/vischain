package core

import (
	"fmt"
	"time"
)

type Blockchain []Block

func NewBlockchain() Blockchain {
	var bc Blockchain
	genesis := Block{
		index:        0,
		timestamp:    time.Now().UnixNano(),
		transactions: nil,
		proof:        nil,
		prevHash:     nil,
	}
	bc = append(bc, genesis)
	return bc
}

func (bc *Blockchain) LastBlock() Block {
	return (*bc)[len(*bc)-1]
}

func (bc *Blockchain) Mine(c chan TransactionList) {
	lastBlock := bc.LastBlock()
	var tl TransactionList
	block := &Block{
		index:     lastBlock.index + 1,
		timestamp: time.Now().UnixNano(),
		prevHash:  lastBlock.Hash(),
	}

	var counter uint32
	counter = 0
	success, hash := ProofOfWork(lastBlock.Hash(), counter)
	for {
		select {
		case tl = <-c:
			block.transactions = tl
		default:
			fmt.Println(success, hash)
			counter++
			success, hash = ProofOfWork(lastBlock.Hash(), counter)
			block.proof = hash
			if success == true {
				fmt.Println("SUCCESS", hash)
				return
			}
		}
	}
}
