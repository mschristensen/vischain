package network

import (
	"encoding/base64"
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
		t := make(map[string]interface{})
		t["code"] = 2
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(t)
		return
	}

	transaction := &core.Transaction{}
	err = transaction.FromMap(m)
	if err != nil { // the transaction could not be parsed
		t := make(map[string]interface{})
		t["code"] = 3
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(t)
		return
	}

	chanT <- *transaction

	t := make(map[string]interface{})
	t["code"] = 1
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

// ReceiveBlock streams out blocks received from peer nodes
func ReceiveBlock(w http.ResponseWriter, r *http.Request, chanB chan core.Block) {
	defer r.Body.Close()

	// parse the request body
	m, err := ParseBody(r.Body)
	if err != nil { // the request body could not be parsed
		t := make(map[string]interface{})
		t["code"] = 2
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(t)
		return
	}

	block := &core.Block{}
	err = block.FromMap(m)
	if err != nil { // the block could not be parsed
		t := make(map[string]interface{})
		t["code"] = 3
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(t)
		return
	}

	chanB <- *block

	t := make(map[string]interface{})
	t["code"] = 1
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

// GetChain
func GetChain(w http.ResponseWriter, r *http.Request, chain core.Blockchain) {
	defer r.Body.Close()

	lastBlockHash := r.URL.Query().Get("lastBlockHash")
	j := 0
	for i := 0; i < len(chain); i++ {
		if base64.StdEncoding.EncodeToString(chain[i].PrevHash) == lastBlockHash {
			j = i
			break
		}
	}

	var code int8
	if j == 0 {
		code = 2
	} else {
		code = 1
	}
	t := make(map[string]interface{})
	t["code"] = code
	t["payload"] = chain[j:]
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}
