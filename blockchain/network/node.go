package network

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/mschristensen/brocoin/blockchain/core"
)

type Node struct {
	Address string
	Peers   []string
	Chain   core.Blockchain
}

func (node *Node) Start(wg *sync.WaitGroup) {

	fmt.Println("Started Node on "+node.Address, node.Chain)

	// set up channels to communicate with goroutines
	minerChanT := make(chan core.Transaction)   // to forward inbound transactions to miner
	minerChanB := make(chan core.Block)         // to receive mined blocks from miner
	minerChanLB := make(chan core.Block)        // to send updated last block on chain to miner
	networkChanT := make(chan core.Transaction) // to receive inbound transactions from network
	networkChanB := make(chan core.Block)       // to receive inbound blocks from network

	// handle incoming requests from peer nodes, streaming out the data along the channels
	go Listen(node, networkChanT, networkChanB)

	// listen for incoming transactions received from peer nodes and forward to the mining process
	go receiveTransactions(networkChanT, minerChanT)

	// mine blocks
	go core.Mine(minerChanLB, minerChanT, minerChanB)

	// set miner up with the initial last block
	lb := node.Chain.LastBlock()
	minerChanLB <- lb

	var bMine core.Block
	var bPeer core.Block
	for {
		select {
		case bMine = <-minerChanB: // we've just mined a block
			last := node.Chain.LastBlock()
			if bMine.Validate(last) { // validate it
				node.Chain.AddBlock(bMine) // add it to the chain

				// broadcast the block to the network
				r, err := Post("/block?peers="+strings.Join(node.Peers, ","), bMine.ToJSON())
				if err != nil {
					log.Fatal("Request to API resulted in an error")
					panic(err)
				}
				if r.StatusCode != 200 {
					log.Fatal(fmt.Sprintf("Request to API did not succeed, got HTTP %d", r.StatusCode))
				}
				m, err := ParseBody(r.Body)
				if err != nil {
					log.Fatal("API response body could not be parsed")
					panic(err)
				}
				response := Response{}
				err = response.FromMap(m)
				if err != nil {
					log.Fatal("API response body could not be written to Response object")
					panic(err)
				}

			} else {
				// TODO: Notify of rejected/invalid block
				fmt.Println("Mined block rejected as invalid", bMine)
			}
		case bPeer = <-networkChanB: // received a block from a peer
			fmt.Println("RECEIVED FROM PEER", bPeer)
		}
	}
}

func receiveTransactions(networkChanT chan core.Transaction, minerChanT chan core.Transaction) {
	for {
		select {
		case t := <-networkChanT:
			minerChanT <- t
		}
	}
}
