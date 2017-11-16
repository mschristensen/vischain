package network

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mschristensen/brocoin/blockchain/core"
)

const APIUrl = "http://localhost:3001/api/v1"

// Listen to incoming requests from peer nodes
// NOTE: `node` is shared across goroutines, treat it as readonly!
func Listen(node *Node, chanT chan core.Transaction, chanB chan core.Block) {
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

// Post the JSON-encoded string `body` to the endpoint `route`
func Post(route string, body string) (*http.Response, error) {
	buf := bytes.NewBuffer([]byte(body))
	r, err := http.Post(APIUrl+route, "application/json; charset=utf-8", buf)

	// if we have a BadGateway error, remove the offline nodes from our list of peers
	if err == nil && r.StatusCode == 502 {
		// TODO: ...
	}

	return r, err
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
