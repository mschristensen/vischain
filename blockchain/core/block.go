package core

import (
	"encoding/json"
	"fmt"
	"time"
)

type Block struct {
	Index        int64				`json:"index,string"`
	Timestamp    int64				`json:"timestamp,string"`
	Transactions TransactionList	`json:"transactions"`
	Proof        Hash				`json:"proof"`
	PrevHash     Hash				`json:"prevHash"`
}

type BroadcastableBlock struct {
	OriginalSender	string		`json:"originalSender"`
	Recipients		[]string	`json:"recipients"`
	Data 			Block		`json:"data"`
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

// ToBroadcastableJSON generates a JSON for a Post request body to the API
// describing the block to be sent, from and to whom.
func (block *Block) ToBroadcastableJSON(sender string, recipients []string) ([]byte, error) {
	b := BroadcastableBlock{
		OriginalSender:	sender,
		Recipients:		recipients,
		Data:			*block,
	}
	data, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}
	return data, nil
}
