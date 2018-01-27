package network

import (
	"sync"
	"bytes"
	"net/http"
	"github.com/gorilla/mux"
	
	"github.com/mschristensen/vischain/blockchain/core"
	"github.com/mschristensen/vischain/blockchain/util"
)

const APIUrl = "http://localhost:3001/api/v1"
type Node struct {
	Address string
	Peers   []string
	Chain   core.Blockchain
	Logger  util.Logger
}

// BlockPackage wraps a block with other meta information
// needed by the recipient when sending across a channel
type BlockPackage struct {
	Sender string
	Block core.Block
}

func (node *Node) Start(wg *sync.WaitGroup) {
	node.Logger.Infof("Started with initial chain: %v", node.Chain)

	// set up channels to communicate with goroutines
	minerChanT := make(chan core.Transaction)   // to forward inbound transactions to miner
	minerChanB := make(chan core.Block)         // to receive mined blocks from miner
	minerChanLB := make(chan core.Block)        // to send updated last block on chain to miner
	networkChanT := make(chan core.Transaction) // to receive inbound transactions from network
	networkChanB := make(chan BlockPackage)     // to receive inbound blocks from network

	// handle incoming requests from peer nodes, streaming out the data along the channels
	go node.Listen(networkChanT, networkChanB)

	// listen for incoming transactions received from peer nodes and forward to the mining process
	go node.receiveTransactions(networkChanT, minerChanT)

	// mine blocks
	go core.Mine(minerChanLB, minerChanT, minerChanB, node.Logger)

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
				broadcastJSON, err := bMine.ToBroadcastableJSON(node.Address, node.Peers)
				if err != nil {
					panic(err)
				}
				r, err := node.Request("POST", "/block", broadcastJSON)
				if err != nil {
					node.Logger.Fatal("Request to API resulted in an error")
					panic(err)
				}
				if r.StatusCode != 200 {
					node.Logger.Fatalf("Request to API did not succeed, got HTTP %d", r.StatusCode)
				}
				node.Logger.Infof("Broadcasting block to peers: %v", bMine)
				// TODO: act on unsuccessful broadcast to nodes
			} else {
				// TODO: Notify of rejected/invalid block
				node.Logger.Infof("Mined block rejected as invalid: %v", bMine)
			}
		case bPeerPackage = <-networkChanB: // received a block from a peer
			bPeer := bPeerPackage.Block
			last := node.Chain.LastBlock()
			if bPeer.Validate(last) {
				// someone has mined a valid block,
				// add it to the chain and inform the mining process
				node.Logger.Infof("Receives valid block and adds to chain: %v", bPeer)
				node.Chain.AddBlock(bPeer)
				minerChanLB <- node.Chain.LastBlock()
			} else if bPeer.Index > lb.Index+1 {
				// the peer has a longer chain than us...
				node.Logger.Info("ATTEMPT CHAIN FETCH")
				r, err := node.Request("GET", "/chain?peer="+bPeerPackage.Sender, nil)
				if err != nil {
					node.Logger.Fatal("Request to API resulted in an error")
					panic(err)
				}
				if r.StatusCode != 200 {
					node.Logger.Fatalf("Request to API did not succeed, got HTTP %d", r.StatusCode)
				}
				// m, err := ParseBody(r.Body)
				// if err != nil {
				// 	log.Fatal("API response body could not be parsed")
				// 	panic(err)
				// }
				// response := Response{}
				// err = response.FromMap(m)
				// if err != nil {
				// 	log.Fatal("API response body could not be written to Response object")
				// 	panic(err)
				// }
				// fmt.Println("CHAIN", response.Payload["payload"])
				// TODO parse chain and update local chain, send new last block to miner
			} else {
				// TODO notify of ignored chain
				node.Logger.Infof("RECEIVED BLOCK INVALID %v %v", bPeer, last)
			}
		}
	}
}

// Request makes an HTTP request
func (node *Node) Request(method string, route string, body []byte) (*http.Response, error) {
	var req *http.Request
	var err error
	if method == "POST" && body != nil {
		req, err = http.NewRequest(method, APIUrl+route, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(method, APIUrl+route, nil)
	}

	if err != nil {
		node.Logger.Fatalf("Error creating HTTP request: %v", err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("X-Sender", node.Address)
	client := &http.Client{}
	r, err := client.Do(req)
	if err != nil {
		node.Logger.Fatalf("Error performing HTTP request: %v", err);
		return nil, err
	}

	return r, nil
}

// Listen to incoming requests from peer nodes
func (node *Node) Listen(chanT chan core.Transaction, chanB chan BlockPackage) {
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/transaction", func(w http.ResponseWriter, r *http.Request) {
		ReceiveTransaction(w, r, chanT)
	}).Methods("POST")
	router.HandleFunc("/block", func(w http.ResponseWriter, r *http.Request) {
		ReceiveBlock(w, r, chanB)
	}).Methods("POST")
	router.HandleFunc("/chain", func(w http.ResponseWriter, r *http.Request) {
		GetChain(w, r, node.Chain)
	}).Methods("GET")

	// Start the server
	http.ListenAndServe(":"+node.Address, router)
}

func (node *Node) receiveTransactions(networkChanT chan core.Transaction, minerChanT chan core.Transaction) {
	for {
		select {
		case t := <-networkChanT:
			minerChanT <- t
		}
	}
}
