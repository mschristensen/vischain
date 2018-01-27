package network

import (
	"net/http"
	"encoding/json"
	"encoding/base64"
	"io/ioutil"
	"fmt"

	"github.com/mschristensen/vischain/blockchain/core"
)

// ReceiveTransaction streams out transactions received from peer nodes
func ReceiveTransaction(w http.ResponseWriter, r *http.Request, chanT chan core.Transaction) {
	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t := make(map[string]interface{})
		t["code"] = 2
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(t)
		return
	}

	var transaction core.Transaction
	err = json.Unmarshal(data, &transaction)
	if err != nil { // the transaction could not be parsed
		fmt.Println("Cannot parse incoming transaction:", err)
		t := make(map[string]interface{})
		t["code"] = 3
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(t)
		return
	}
	fmt.Println(transaction)

	chanT <- transaction

	t := make(map[string]interface{})
	t["code"] = 1
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

// ReceiveBlock streams out blocks received from peer nodes
func ReceiveBlock(w http.ResponseWriter, r *http.Request, chanB chan BlockPackage) {
	defer r.Body.Close()

	// parse the request body
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t := make(map[string]interface{})
		t["code"] = 2
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(t)
		return
	}

	var block core.Block
	err = json.Unmarshal(data, &block)
	if err != nil { // the block could not be parsed
		fmt.Println("Cannot parse incoming block:", err)
		t := make(map[string]interface{})
		t["code"] = 3
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(t)
		return
	}

	blockPackage := &BlockPackage{
		Sender: r.Header.Get("X-Sender"),
		Block: block,
	}
	chanB <- *blockPackage

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
