package main

import (
	"fmt"

	"github.com/mschristensen/brocoin/blockchain"
)

func printBlock(b blockchain.Block) {
	fmt.Println(b.PrevHash)
}

func main() {
	fmt.Println("Hello, world!")
	bc := blockchain.Init()
	printBlock(bc[0])
}
