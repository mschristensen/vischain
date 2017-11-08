package main

import (
	"fmt"

	"github.com/mschristensen/brocoin/blockchain/api"
	"github.com/mschristensen/brocoin/blockchain/core"
)

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

	r := api.HelloResponseGet{}
	api.Get("/hello", &r)
	fmt.Println(r)

	r2 := api.HelloResponseGet{
		Payload: "Test Payload",
		Title:   "Test Title",
	}
	r3 := api.HelloResponsePost{}
	api.Post("/hello", r2, &r3)
	fmt.Println(r3)

	api.Listen()
}

func randomTransaction() core.Transaction {
	return core.Transaction{
		Sender:    "from",
		Recipient: "to",
		Amount:    1,
	}
}
