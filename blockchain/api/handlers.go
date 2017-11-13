package api

import (
	"encoding/json"
	"net/http"

	"github.com/mschristensen/brocoin/blockchain/core"
)

// ReceiveTransaction streams out transactions received from peer nodes
func ReceiveTransaction(w http.ResponseWriter, r *http.Request, chanT chan core.Transaction) {
	defer r.Body.Close()

	// parse the request body
	m, err := ParseBody(r.Body)
	if err != nil { // the request body could not be parsed
		t := PeerResponse{Code: 2}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(t)
		return
	}

	transaction := &core.Transaction{}
	err = transaction.FromMap(m)
	if err != nil { // the transaction could not be parsed
		t := PeerResponse{Code: 3}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(t)
		return
	}

	chanT <- *transaction

	t := PeerResponse{Code: 1}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

// ReceiveBlock streams out blocks received from peer nodes
func ReceiveBlock(w http.ResponseWriter, r *http.Request, chanB chan core.Block) {
	defer r.Body.Close()

	// parse the request body
	m, err := ParseBody(r.Body)
	if err != nil { // the request body could not be parsed
		t := PeerResponse{Code: 2}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(t)
		return
	}

	block := &core.Block{}
	err = block.FromMap(m)
	if err != nil { // the block could not be parsed
		t := PeerResponse{Code: 3}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(t)
		return
	}

	chanB <- *block

	t := PeerResponse{Code: 1}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}
