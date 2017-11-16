package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/mschristensen/brocoin/blockchain/core"
	"github.com/mschristensen/brocoin/blockchain/network"
)

func main() {
	// read network configuration from file
	inFile, _ := os.Open("network.config")
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	// create the nodes
	var nodes []network.Node
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		nodes = append(nodes, network.Node{
			Address: parts[0],
			Peers:   parts[1:],
			Chain:   core.NewBlockchain(),
		})
	}

	// start the nodes
	var wg sync.WaitGroup
	for i := 0; i < len(nodes); i++ {
		fmt.Println(nodes[i].Address)
		wg.Add(1)
		go nodes[i].Start(&wg)
	}

	// Wait for all HTTP fetches to complete.
	wg.Wait()
}
