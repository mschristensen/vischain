package core

import (
	"encoding/binary"
	"time"

	"github.com/mschristensen/vischain/blockchain/util"
)

type Blockchain []Block

type BroadcastableBlockchain struct {
	Code	uint16		`json:"code"`
	Chain	Blockchain	`json:"payload"`
}

func NewBlockchain() Blockchain {
	var bc Blockchain
	genesis := Block{
		Index:        0,
		Timestamp:    time.Now().UnixNano(),
		Transactions: nil,
		Proof:        nil,
		PrevHash:     nil,
	}
	bc = append(bc, genesis)
	return bc
}

func (bc *Blockchain) LastBlock() Block {
	return (*bc)[len(*bc)-1]
}

// AddBlock adds a block to the blockchain
func (bc *Blockchain) AddBlock(block Block) {
	*bc = append(*bc, block)
}

// ProofOfWork implements a HashCash-based PoW algorithm.
// Work is completed successfully if the hash of the block and a
// given counter ends with a sufficient number of trailing zeroes.
func ProofOfWork(blockHash []byte, counter []byte) bool {
	bs := util.ConcatBytes(blockHash, counter)
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

// Validate indicates whether an entire blockchain is valid
func (bc *Blockchain) Validate() bool {
	for i := 0; i < len(*bc); i++ {
		if i == 0 {
			continue
		}
		if !(*bc)[i].Validate((*bc)[i-1]) {
			return false
		}
	}
	return true
}

// Mine continuously accepts new transactions and attempts to mine a block containing
// them by finding a proof value which satisfies the difficulty constraint
func Mine(chanLB chan Block, chanT chan Transaction, chanB chan Block, log util.Logger) {
	var t Transaction // incoming transaction
	var lb Block      // current last block on the chain
	var block Block   // block to mine

	var counterInt32 uint32
	counter := []byte{0, 0, 0, 0}
	success := false
	for {
		select {
		case t = <-chanT:
			block.Transactions.AddTransaction(t)
			log.Infof("MINER RECEIVED TRANSACTION: %v", t)
		case lb = <-chanLB:
			block = lb.NewBlock()
			log.Infof("MINER RECEIVED BLOCK: %v", lb)
		default:
			if block.Transactions == nil || &lb == nil {
				continue
			}
			// increment counter
			counterInt32 = binary.LittleEndian.Uint32(counter)
			counterInt32++
			binary.LittleEndian.PutUint32(counter, counterInt32)

			success = ProofOfWork(lb.Hash(), counter)
			if success == true {
				block.Proof = append([]byte(nil), counter...)
				chanB <- block
				log.Infof("MINED NEW BLOCK: %v", block)
				lb = block
				block = lb.NewBlock()
			}
		}
	}
}
