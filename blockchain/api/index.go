package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

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
