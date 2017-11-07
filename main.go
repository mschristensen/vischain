package main

import (
	"fmt"

	"github.com/mschristensen/brocoin/blockchain/api"
	"github.com/mschristensen/brocoin/blockchain/core"
)

type Response struct {
	Payload string
	Title   string
}

func main() {
	bc := core.NewBlockchain()

	c := make(chan core.Transaction)
	go bc.Mine(c)

	// Add transactions until limit reached
	for i := 0; i < 1000; i++ {
		c <- randomTransaction()
	}

	fmt.Println(bc)
	fmt.Println(bc.ValidateBlockchain())

	var resp Response
	api.Request("/hello", &resp)
	fmt.Println(resp)

	api.Listen()
}

func randomTransaction() core.Transaction {
	return core.Transaction{
		Sender:    "from",
		Recipient: "to",
		Amount:    1,
	}
}
