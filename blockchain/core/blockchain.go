package core

import (
	"encoding/binary"
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

func (lb *Block) NewBlock() Block {
	return Block{
		index:        lb.index + 1,
		timestamp:    time.Now().UnixNano(),
		transactions: nil,
		proof:        nil,
		prevHash:     lb.Hash(),
	}
}

// AddBlock adds a block to the blockchain
func (bc *Blockchain) AddBlock(block Block) {
	*bc = append(*bc, block)
}

// ValidateBlock indicates whether a new block is valid, which it is iff.
// its proof hashed with the last block hash satisfies the difficulty constraint
func (bc *Blockchain) ValidateBlock(lastBlock Block, block Block) bool {
	return CompareHashes(block.prevHash, lastBlock.Hash()) && ProofOfWork(block.prevHash, block.proof)
}

// ProofOfWork implements a HashCash-based PoW algorithm.
// Work is completed successfully if the hash of the block and a
// given counter ends with a sufficient number of trailing zeroes.
func ProofOfWork(blockHash []byte, counter []byte) bool {
	bs := ConcatBytes(blockHash, counter)
	difficulty := 2
	hash := Sha256(bs)
	tail := hash[len(hash)-difficulty:]
	for i := 0; i < difficulty; i++ {
		if tail[i] != 0 {
			return false
		}
	}

	return true
}

// ValidateBlockchain indicates whether an entire blockchain is valid
func (bc *Blockchain) ValidateBlockchain() bool {
	for i := 0; i < len(*bc); i++ {
		if i == 0 {
			continue
		}
		if !bc.ValidateBlock((*bc)[i-1], (*bc)[i]) {
			return false
		}
	}
	return true
}

// Mine continuously accepts new transactions and attempts to mine a block containing
// them by finding a proof value which satisfies the difficulty constraint
func Mine(chanLB chan *Block, chanT chan *Transaction, chanB chan *Block) {
	var t *Transaction // incoming transaction
	var lb *Block      // current last block on the chain
	var block Block    // block to mine

	var counterInt32 uint32
	counter := []byte{0, 0, 0, 0}
	// lastBlock := bc.LastBlock()
	// success := ProofOfWork(lb.Hash(), counter)
	success := false
	for {
		select {
		case t = <-chanT:
			block.transactions.AddTransaction(*t)
			fmt.Println("RECEIVED TRANSACTION", *t)
		case lb = <-chanLB:
			block = lb.NewBlock()
			fmt.Println("RECEIVED BLOCK", *lb)
		default:
			if block.transactions == nil || lb == nil {
				continue
			}
			// increment counter
			counterInt32 = binary.LittleEndian.Uint32(counter)
			counterInt32++
			binary.LittleEndian.PutUint32(counter, counterInt32)

			success = ProofOfWork(lb.Hash(), counter)
			if success == true {
				block.proof = append([]byte(nil), counter...)
				// bc.AddBlock(block)
				chanB <- &block
				fmt.Println("MINED", block)
				// TODO broadcast block to peers
				lb = &block
				block = lb.NewBlock()
			}
		}
	}
}
