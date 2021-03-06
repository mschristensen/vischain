package main

import (
	"bufio"
	"os"
	"strings"
	"sync"
	"context"

	"github.com/mschristensen/vischain/blockchain/core"
	"github.com/mschristensen/vischain/blockchain/network"
	"github.com/mschristensen/vischain/blockchain/util"
)

func main() {
	// read network configuration from file
	inFile, _ := os.Open("network.config")
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	// create the nodes, with identical initial blockchains
	var nodes []network.Node
	chain := core.NewBlockchain()
	i := 0
	ctxBackground := context.Background()
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		ctx := util.CreateLogger(ctxBackground, parts[0])
		chainCopy := core.Blockchain{}
		nodes = append(nodes, network.Node{
			Address: parts[0],
			Peers:   parts[1:],
			Chain:   append(chainCopy, chain...),
			Logger:	 util.GetLogger(ctx),
		})
		i++
	}

	// start the nodes
	var wg sync.WaitGroup
	for i := 0; i < len(nodes); i++ {
		wg.Add(1)
		go nodes[i].Start(&wg)
	}

	wg.Wait()
}
