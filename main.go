package main

import (
	"fmt"

	"github.com/mschristensen/brocoin/core"
)

func main() {
	bc := core.NewBlockchain()

	// Add transactions until limit reached
	var tl core.TransactionList
	t := randomTransaction()
	for tl.AddTransaction(t) == 1 {
		t = randomTransaction()
	}

	fmt.Println(bc)
}

func randomTransaction() core.Transaction {
	return core.Transaction{
		Sender:    "from",
		Recipient: "to",
		Amount:    1,
	}
}
