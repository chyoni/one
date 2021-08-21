package main

import (
	"fmt"

	"github.com/chiwon99881/one/blockchain"
)

func main() {
	chain := blockchain.BlockChain()
	for _, block := range chain.Blocks {
		fmt.Println(block)
	}
}
