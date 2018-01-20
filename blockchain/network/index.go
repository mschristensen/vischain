package network

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"log"

	"github.com/gorilla/mux"
	"github.com/mschristensen/vischain/blockchain/core"
)

const APIUrl = "http://localhost:3001/api/v1"

// Listen to incoming requests from peer nodes
// NOTE: `node` is shared across goroutines, treat it as readonly!
func Listen(node *Node, chanT chan core.Transaction, chanB chan BlockPackage) {
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

// Request makes an HTTP request
func Request(method string, route string, body string, sender string) (*http.Response, error) {
	var req *http.Request
	var err error
	if method == "POST" && body != "" {
		req, err = http.NewRequest(method, APIUrl+route, bytes.NewBuffer([]byte(body)))
	} else {
		req, err = http.NewRequest(method, APIUrl+route, nil)
	}

	if err != nil {
		log.Fatal("Error creating HTTP request: ", err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("X-Sender", sender)
	client := &http.Client{}
	r, err := client.Do(req)
	if err != nil {
		log.Fatal("Error performing HTTP request: ", err);
		return nil, err
	}

	return r, nil
}

// ParseBody parses application/json data of unknown shape
// from a response body into an empty interface map
func ParseBody(body io.Reader) (map[string]interface{}, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var parsedMap map[string]interface{}
	err = json.Unmarshal(data, &parsedMap)
	if err != nil {
		return nil, err
	}
	return parsedMap, nil
}
