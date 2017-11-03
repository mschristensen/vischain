package core

import "fmt"

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
