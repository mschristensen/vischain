package main

import (
	"fmt"

	"github.com/mschristensen/brocoin/core"
)

func main() {
	bc := core.NewBlockchain()
	(&bc).AddBlock(1, "abc")
	fmt.Println(bc)
}
