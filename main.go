package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mschristensen/brocoin/blockchain/api"
	"github.com/mschristensen/brocoin/blockchain/core"
	"github.com/mschristensen/brocoin/blockchain/node"
)

func main() {
	// init this node
	args := os.Args[1:]
	self := node.Init(args[0], args[1:])

	bc := core.NewBlockchain()

	minerChanT := make(chan core.Transaction) // to forward inbound transactions to miner
	minerChanB := make(chan core.Block)       // to receive mined blocks from miner
	minerChanLB := make(chan core.Block)      // to send updated last block on chain to miner
	apiChanT := make(chan core.Transaction)   // to receive inbound transactions from api
	apiChanB := make(chan core.Block)         // to receive inbound blocks from api

	go api.Listen(self.Address, apiChanT, apiChanB)
	go core.Mine(minerChanLB, minerChanT, minerChanB)
	go receiveTransactions(apiChanT, minerChanT)

	// set miner up with the initial last block
	lb := bc.LastBlock()
	minerChanLB <- lb

	var bMine core.Block
	var bPeer core.Block
	for {
		select {
		case bMine = <-minerChanB: // we've just mined a block.
			last := bc.LastBlock()
			if bMine.Validate(last) { // validate it
				bc.AddBlock(bMine) // add it to the chain

				// broadcast the block to the network
				// TODO: handle errors
				r, _ := api.Post("/block?peers="+strings.Join(self.Peers, ","), bMine.ToJSON())
				m, _ := api.ParseBody(r.Body)
				response := api.APIResponse{}
				response.FromMap(m)
			} else {
				// TODO: Notify of rejected/invalid block
				fmt.Println("Mined block rejected as invalid", bMine)
			}
		case bPeer = <-apiChanB: // received a block from a peer
			fmt.Println("RECEIVED FROM PEER", bPeer)
		}
	}

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

	// Strategy
	//      Create a worker which listens to events and dispatches to relevant farmers

	// Events:
	//      Listen for transactions from peers
	//      Listen for invoked transactions
	//      Listen for blocks from peers
	//      Mine continuously + emit blocks
}

func receiveTransactions(apiChanT chan core.Transaction, minerChanT chan core.Transaction) {
	for {
		select {
		case t := <-apiChanT:
			minerChanT <- t
		}
	}
}
