package api

import (
	"io"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}
