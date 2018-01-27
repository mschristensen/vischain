package core

import (
)

const maxTransactions = 5

type Transaction struct {
	Sender    string	`json:"sender"`
	Recipient string	`json:"recipient"`
	Amount    uint16	`json:"amount,string"`
}

type TransactionList []Transaction

// AddTransaction adds a given transaction to the receiving TransactionList.
// Returns -1 if the TransactionList is already larger than allowed and so the transaction was not added.
// Returns 0  if the transaction was added and the TransactionList is now at its max allowable length.
// Returns 1  if the transaction was added and further transactions may also be added to the TransactionList.
func (tl *TransactionList) AddTransaction(transaction Transaction) int8 {
	if len(*tl) >= maxTransactions {
		return -1
	}
	*tl = append(*tl, transaction)
	if len(*tl) == maxTransactions {
		return 0
	}
	return 1
}