package core

import (
	"fmt"
	"strconv"
)

const maxTransactions = 5

type Transaction struct {
	Sender    string
	Recipient string
	Amount    uint16
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

// ToJSON encodes the transaction as a JSON string.
// The values are represented as:
//      sender          string
//      recipient       string
//      amount          decimal string
func (t *Transaction) ToJSON() string {
	return fmt.Sprintf(`{
        "sender": "%v",
        "recipient": "%v",
        "amount": "%d"
    }`, t.Sender, t.Recipient, t.Amount)
}

// FromMap accepts an empty interface slice describing a transaction list
// (parsed from a JSON) and stores its parsed contents in the TransactionList.
func (tl *TransactionList) FromMap(m []interface{}) {
	var transactions TransactionList
	for _, t := range m {
		transaction := &Transaction{}
		transaction.FromMap(t.(map[string]interface{}))
		transactions = append(transactions, *transaction)
	}
	*tl = transactions
}

// FromMap accepts an empty interface map describing a transaction
// (parsed from a JSON) and stores its parsed contents in the Transaction.
func (t *Transaction) FromMap(m map[string]interface{}) {
	amount, _ := strconv.ParseUint(m["amount"].(string), 10, 16)

	*t = Transaction{
		Sender:    m["sender"].(string),
		Recipient: m["recipient"].(string),
		Amount:    uint16(amount),
	}
}
