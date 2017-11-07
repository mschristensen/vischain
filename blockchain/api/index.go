package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

const APIUrl = "http://localhost:3001/api/v1"

// TODO:
//      AddBlock
//          New block received from peer
//          POST /block
//

// Listen to incoming requests from peer nodes
func Listen() {
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/hello", Hello)

	// Start the server
	http.ListenAndServe(":8080", router)
}

func Request(route string, target interface{}) error {
	resp, err := http.Get(APIUrl + route)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}
