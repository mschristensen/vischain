package core

import (
	"bytes"
	"encoding/binary"
)

type Block struct {
	index        int64
	timestamp    int64
	transactions TransactionList
	proof        Proof
	prevHash     string
}

func (block *Block) Hash() string {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, *block)
	return Sha256(buf.Bytes())
}
