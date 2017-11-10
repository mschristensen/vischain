package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mschristensen/brocoin/blockchain/core"
)

func ReceiveTransaction(w http.ResponseWriter, r *http.Request, chanT chan core.Transaction) {
	defer r.Body.Close()
	m, _ := ParseBody(r.Body)

	transaction := &core.Transaction{}
	transaction.FromMap(m)
	fmt.Println("RECEIVED", transaction)

	chanT <- *transaction

	t := OKResponse{Code: 1}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func ReceiveBlock(w http.ResponseWriter, r *http.Request, chanB chan core.Block) {
	defer r.Body.Close()
	m, _ := ParseBody(r.Body)

	block := &core.Block{}
	block.FromMap(m)
	fmt.Println("RECEIVED", block)

	chanB <- *block

	t := OKResponse{Code: 1}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}
