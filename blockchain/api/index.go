package api

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

// TODO:
//      AddBlock
//          New block received from peer
//          POST /block
//

// Listen to incoming requests from peer nodes
func Listen(c_t chan *core.Transaction) {
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/hello", Hello)
	router.HandleFunc("/transaction", func(w http.ResponseWriter, r *http.Request) {
		ReceiveTransaction(w, r, c_t)
	}).Methods("POST")

	// Start the server
	http.ListenAndServe(":8080", router)
}

func Get(route string, target interface{}) error {
	resp, err := http.Get(APIUrl + route)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

// Post the JSON-encoded string `body` to the endpoint `route`
func Post(route string, body string) (*http.Response, error) {
	buf := bytes.NewBuffer([]byte(body))
	return http.Post(APIUrl+route, "application/json; charset=utf-8", buf)
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
