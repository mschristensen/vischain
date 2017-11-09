package core

import (
	"encoding/binary"
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

func (bc *Blockchain) NewBlock() Block {
	lastBlock := bc.LastBlock()
	return Block{
		index:        lastBlock.index + 1,
		timestamp:    time.Now().UnixNano(),
		transactions: nil,
		proof:        nil,
		prevHash:     lastBlock.Hash(),
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
func (bc *Blockchain) Mine(c chan Transaction) {
	var t Transaction
	block := bc.NewBlock()

	var counterInt32 uint32
	counter := []byte{0, 0, 0, 0}
	lastBlock := bc.LastBlock()
	success := ProofOfWork(lastBlock.Hash(), counter)
	for {
		select {
		case t = <-c:
			block.transactions.AddTransaction(t)
		default:
			if block.transactions == nil {
				continue
			}
			// increment counter
			counterInt32 = binary.LittleEndian.Uint32(counter)
			counterInt32++
			binary.LittleEndian.PutUint32(counter, counterInt32)

			success = ProofOfWork(lastBlock.Hash(), counter)
			if success == true {
				block.proof = append([]byte(nil), counter...)
				bc.AddBlock(block)
				// TODO broadcast block to peers
				lastBlock = block
				block = bc.NewBlock()
			}
		}
	}
}
