package main

import (
	"fmt"

	"github.com/mschristensen/brocoin/core"
)

func main() {
	bc := core.NewBlockchain()
	(&bc).AddBlock(1, "abc", nil)
	t := core.Transaction{
		Sender:    "mike",
		Recipient: "james",
		Amount:    1,
	}
	bc[1].AddTransaction(t)
	fmt.Println(bc)
}
