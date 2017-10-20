package core

type Block struct {
	index        int64
	timestamp    int64
	transactions []Transaction
	proof        uint32
	PrevHash     string
}

func (block *Block) AddTransaction(transaction Transaction) {
	block.transactions = append(block.transactions, transaction)
}
