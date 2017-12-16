package core

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
)

type Block struct {
	Index        int64
	Timestamp    int64
	Transactions TransactionList
	Proof        Hash
	PrevHash     Hash
}

func (block *Block) Hash() Hash {
	return Sha256([]byte(fmt.Sprintf("%v", *block)))
}

func (lb *Block) NewBlock() Block {
	return Block{
		Index:        lb.Index + 1,
		Timestamp:    time.Now().UnixNano(),
		Transactions: nil,
		Proof:        nil,
		PrevHash:     lb.Hash(),
	}
}

// Validate indicates whether a new block is valid, which it is iff.
// its proof hashed with the last block hash satisfies the difficulty constraint
func (block *Block) Validate(lb Block) bool {
	return CompareHashes(block.PrevHash, lb.Hash()) && ProofOfWork(block.PrevHash, block.Proof)
}

// ToJSON encodes the block as a JSON string.
// The values are represented as:
//      index           decimal string
//      timestamp       decimal string
//      transactions    JSON representing transactions
//      proof           base64 encoded string
//      PrevHash        PrevHash encoded string
func (block *Block) ToJSON() string {
	transactionsJSON := "["
	for i, transaction := range block.Transactions {
		transactionsJSON += transaction.ToJSON()

		if i < len(block.Transactions)-1 {
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
    }`, block.Index, block.Timestamp, transactionsJSON, base64.StdEncoding.EncodeToString(block.Proof), base64.StdEncoding.EncodeToString(block.PrevHash))
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
		Index:        index,
		Timestamp:    timestamp,
		Transactions: transactions,
		Proof:        proof,
		PrevHash:     prevHash,
	}
	return nil
}

// ToAPIJSON generates a JSON for a Post request body to the API
// describing the block to be sent, from and to whom.
func (block *Block) ToAPIJSON(sender string, recipients []string) string {
	recipientsJSON := "["
	for i, recipient := range recipients {
		recipientsJSON += (fmt.Sprintf("\"%v\"", recipient))
		if i != len(recipients)-1 {
			recipientsJSON += ","
		}
	}
	recipientsJSON += "]"
	return fmt.Sprintf(`{
		"originalSender": "%v",
		"recipients": %v,
		"data": %v
	}`, sender, recipientsJSON, block.ToJSON())
}
