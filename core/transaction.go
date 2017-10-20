package core

type Transaction struct {
	Sender    string
	Recipient string
	Amount    uint16
}

type TransactionList []Transaction

func (tl *TransactionList) AddTransaction(transaction Transaction) {
	*tl = append(*tl, transaction)
}
