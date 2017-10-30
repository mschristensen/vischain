package main

import (
	"fmt"

	"github.com/mschristensen/brocoin/core"
)

func main() {
	bc := core.NewBlockchain()

	c := make(chan core.Transaction)
	go bc.Mine(c)

	// Add transactions until limit reached
	for i := 0; i < 20; i++ {
		c <- randomTransaction()
	}

	fmt.Println("BLOCKCHAIN", bc)
}

func randomTransaction() core.Transaction {
	return core.Transaction{
		Sender:    "from",
		Recipient: "to",
		Amount:    1,
	}
}
