package network

import (
	"fmt"
	"log"
	"sync"

	"github.com/mschristensen/vischain/blockchain/core"
)

type Node struct {
	Address string
	Peers   []string
	Chain   core.Blockchain
}

// BlockPackage wraps a block with other meta information
// needed by the recipient when sending across a channel
type BlockPackage struct {
	Sender string
	Block core.Block
}

func (node *Node) Start(wg *sync.WaitGroup) {

	fmt.Println("Started Node on "+node.Address, node.Chain)

	// set up channels to communicate with goroutines
	minerChanT := make(chan core.Transaction)   // to forward inbound transactions to miner
	minerChanB := make(chan core.Block)         // to receive mined blocks from miner
	minerChanLB := make(chan core.Block)        // to send updated last block on chain to miner
	networkChanT := make(chan core.Transaction) // to receive inbound transactions from network
	networkChanB := make(chan BlockPackage)     // to receive inbound blocks from network

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
	var bPeerPackage BlockPackage
	for {
		select {
		case bMine = <-minerChanB: // we've just mined a block
			last := node.Chain.LastBlock()
			if bMine.Validate(last) { // validate it
				node.Chain.AddBlock(bMine) // add it to the chain

				// broadcast the block to the network
				r, err := Request("POST", "/block", bMine.ToAPIJSON(node.Address, node.Peers), node.Address)
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
		case bPeerPackage = <-networkChanB: // received a block from a peer
			bPeer := bPeerPackage.Block
			last := node.Chain.LastBlock()
			if bPeer.Validate(last) {
				// someone has mined a valid block,
				// add it to the chain and inform the mining process
				node.Chain.AddBlock(bPeer)
				minerChanLB <- node.Chain.LastBlock()
			} else if bPeer.Index > lb.Index+1 {
				// the peer has a longer chain than us...
				r, err := Request("GET", "/chain?peer="+bPeerPackage.Sender, "", node.Address)
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
				fmt.Println("CHAIN", response.Payload["payload"])
				// TODO parse chain and update local chain, send new last block to miner

			} else {
				// TODO notify of ignored chain
			}
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
