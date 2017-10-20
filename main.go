package main

import (
	"fmt"

	"github.com/mschristensen/brocoin/core"
)

func main() {
	bc := core.NewBlockchain()

	var tl core.TransactionList
	t1 := core.Transaction{
		Sender:    "mike",
		Recipient: "james",
		Amount:    1,
	}
	t2 := core.Transaction{
		Sender:    "bob",
		Recipient: "harry",
		Amount:    5,
	}
	tl.AddTransaction(t1)
	tl.AddTransaction(t2)

	(&bc).AddBlock(1, "abc", tl)
	fmt.Println(bc)
}
