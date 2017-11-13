package main

import (
	"fmt"
	"log"
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

	// init the blockchain
	bc := core.NewBlockchain()

	// set up channels to communicate with goroutines
	minerChanT := make(chan core.Transaction) // to forward inbound transactions to miner
	minerChanB := make(chan core.Block)       // to receive mined blocks from miner
	minerChanLB := make(chan core.Block)      // to send updated last block on chain to miner
	apiChanT := make(chan core.Transaction)   // to receive inbound transactions from api
	apiChanB := make(chan core.Block)         // to receive inbound blocks from api

	// handle incoming requests from peer nodes, streaming out the data along the channels
	go api.Listen(self.Address, apiChanT, apiChanB)

	// listen for incoming transactions received from peer nodes and forward to the mining process
	go receiveTransactions(apiChanT, minerChanT)

	// mine blocks
	go core.Mine(minerChanLB, minerChanT, minerChanB)

	// set miner up with the initial last block
	lb := bc.LastBlock()
	minerChanLB <- lb

	var bMine core.Block
	var bPeer core.Block
	for {
		select {
		case bMine = <-minerChanB: // we've just mined a block
			last := bc.LastBlock()
			if bMine.Validate(last) { // validate it
				bc.AddBlock(bMine) // add it to the chain

				// broadcast the block to the network
				r, err := api.Post("/block?peers="+strings.Join(self.Peers, ","), bMine.ToJSON())
				if err != nil {
					log.Fatal("Request to API resulted in an error")
					panic(err)
				}
				if r.StatusCode != 200 {
					log.Fatal(fmt.Sprintf("Request to API did not succeed, got HTTP %d", r.StatusCode))
				}
				m, err := api.ParseBody(r.Body)
				if err != nil {
					log.Fatal("API response body could not be parsed")
					panic(err)
				}
				response := api.Response{}
				err = response.FromMap(m)
				if err != nil {
					log.Fatal("API response body could not be written to Response object")
					panic(err)
				}

			} else {
				// TODO: Notify of rejected/invalid block
				fmt.Println("Mined block rejected as invalid", bMine)
			}
		case bPeer = <-apiChanB: // received a block from a peer
			fmt.Println("RECEIVED FROM PEER", bPeer)
		}
	}
}

func receiveTransactions(apiChanT chan core.Transaction, minerChanT chan core.Transaction) {
	for {
		select {
		case t := <-apiChanT:
			minerChanT <- t
		}
	}
}
