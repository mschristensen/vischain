package main

import (
	"fmt"

	"github.com/mschristensen/brocoin/blockchain"
)

func main() {
	bc := blockchain.Init()
	(&bc).AddBlock(1, "abc")
	fmt.Println(bc)
}
