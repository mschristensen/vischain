package core

import (
	"bytes"
	"encoding/binary"
)

type Block struct {
	index        int64
	timestamp    int64
	transactions TransactionList
	proof        Hash
	prevHash     Hash
}

func (block Block) Hash() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, block)
	return Sha256(buf.Bytes())
}
