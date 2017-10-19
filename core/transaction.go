package core

type Transaction struct {
	Sender    string
	Recipient string
	Amount    uint16
}

func (block *Block) AddTransaction(transaction Transaction) {
	block.transactions = append(block.transactions, transaction)
}
