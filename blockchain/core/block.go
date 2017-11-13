package core

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
)

type Block struct {
	index        int64
	timestamp    int64
	transactions TransactionList
	proof        Hash
	prevHash     Hash
}

func (block *Block) Hash() Hash {
	return Sha256([]byte(fmt.Sprintf("%v", *block)))
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

// Validate indicates whether a new block is valid, which it is iff.
// its proof hashed with the last block hash satisfies the difficulty constraint
func (block *Block) Validate(lb Block) bool {
	return CompareHashes(block.prevHash, lb.Hash()) && ProofOfWork(block.prevHash, block.proof)
}

// ToJSON encodes the block as a JSON string.
// The values are represented as:
//      index           decimal string
//      timestamp       decimal string
//      transactions    JSON representing transactions
//      proof           base64 encoded string
//      prevHash        prevHash encoded string
func (block *Block) ToJSON() string {
	transactionsJSON := "["
	for i, transaction := range block.transactions {
		transactionsJSON += transaction.ToJSON()

		if i < len(block.transactions)-1 {
			transactionsJSON += ","
		}
	}
	transactionsJSON += "]"

	return fmt.Sprintf(`{
        "index": "%d",
        "timestamp": "%d",
        "transactions": %v,
        "proof": "%v",
        "prevHash": "%v"
    }`, block.index, block.timestamp, transactionsJSON, base64.StdEncoding.EncodeToString(block.proof), base64.StdEncoding.EncodeToString(block.prevHash))
}

// FromMap accepts an empty interface map describing a block
// (parsed from a JSON) and stores its parsed contents in the Block.
func (block *Block) FromMap(m map[string]interface{}) error {
	index, err := strconv.ParseInt(m["index"].(string), 10, 64)
	if err != nil {
		return err
	}
	timestamp, err := strconv.ParseInt(m["timestamp"].(string), 10, 64)
	if err != nil {
		return err
	}
	transactions := TransactionList{}
	err = transactions.FromMap(m["transactions"].([]interface{}))
	if err != nil {
		return err
	}
	proof, err := base64.StdEncoding.DecodeString(m["proof"].(string))
	if err != nil {
		return err
	}
	prevHash, err := base64.StdEncoding.DecodeString(m["prevHash"].(string))
	if err != nil {
		return err
	}

	*block = Block{
		index:        index,
		timestamp:    timestamp,
		transactions: transactions,
		proof:        proof,
		prevHash:     prevHash,
	}
	return nil
}
