package main

import (
	"fmt"
	"os"

	"github.com/mschristensen/brocoin/blockchain/api"
	"github.com/mschristensen/brocoin/blockchain/core"
	"github.com/mschristensen/brocoin/blockchain/node"
)

func main() {
	// init this node
	args := os.Args[1:]
	node.Init(args[0], args[1:])

	bc := core.NewBlockchain()

	minerChanT := make(chan *core.Transaction)
	apiChanT := make(chan *core.Transaction)

	go bc.Mine(minerChanT)

	go receiveTransactions(apiChanT, minerChanT)
	// Add transactions until limit reached
	// for i := 0; i < 1000; i++ {
	// 	c <- randomTransaction()
	// }

	// fmt.Println(bc)
	// fmt.Println(bc.ValidateBlockchain())

	lb := bc.LastBlock()
	res, _ := api.Post("/hello", lb.ToJSON())
	m, _ := api.ParseBody(res.Body)
	block := &core.Block{}
	a := api.APIResponse{}
	a.FromMap(m)
	block.FromMap(a.Payload)
	fmt.Println("MAP", a.Payload)
	fmt.Println("BLOCK SENT    ", lb)
	fmt.Println("BLOCK RECEIVED", *block)

	api.Listen(apiChanT)

	// Strategy
	//      Create a worker which listens to events and dispatches to relevant farmers

	// Events:
	//      Listen for transactions from peers
	//      Listen for invoked transactions
	//      Listen for blocks from peers
	//      Mine continuously + emit blocks
}

func receiveTransactions(apiChanT chan *core.Transaction, minerChanT chan *core.Transaction) {
	for {
		select {
		case t := <-apiChanT:
			minerChanT <- t
		}
	}
}
