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

func (bc *Blockchain) NewBlock() Block {
	return Block{
		index:        bc.LastBlock().index + 1,
		timestamp:    time.Now().UnixNano(),
		prevHash:     bc.LastBlock().Hash(),
		transactions: nil,
	}
}

func (bc *Blockchain) AddBlock(block Block) {
	*bc = append(*bc, block)
}

func (bc *Blockchain) Mine(c chan Transaction) {
	var t Transaction
	block := bc.NewBlock()

	var counterInt32 uint32
	counter := []byte{0, 0, 0, 0}
	success := ProofOfWork(bc.LastBlock().Hash(), counter)
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

			success = ProofOfWork(bc.LastBlock().Hash(), counter)
			if success == true {
				block.proof = counter
				bc.AddBlock(block)
				fmt.Println(bc.LastBlock().transactions)
				// TODO broadcast block to peers
				block = bc.NewBlock()
			}
		}
	}
}
