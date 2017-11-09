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

	minerChanT := make(chan *core.Transaction) // to forward inbound transactions to miner
	minerChanB := make(chan *core.Block)       // to receive mined blocks from miner
	minerChanLB := make(chan *core.Block)      // to send updated last block on chain to miner
	apiChanT := make(chan *core.Transaction)   // to receive inbound transactions from api

	go core.Mine(minerChanLB, minerChanT, minerChanB)
	go receiveMinedBlocks(minerChanLB, minerChanB)
	go receiveTransactions(apiChanT, minerChanT)

	// set miner up with the initial last block
	lb := bc.LastBlock()
	minerChanLB <- &lb

	// fmt.Println(bc)
	// fmt.Println(bc.ValidateBlockchain())

	// lb := bc.LastBlock()
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

// TODO lock bc before forwarding on
func receiveMinedBlocks(chanLB chan *core.Block, chanB chan *core.Block) {
	for {
		select {
		case b := <-chanB:
			bc.AddBlock(*b)
		}
	}
}
