package main

import (
	"github.com/chiwon99881/one/blockchain"
	"github.com/chiwon99881/one/db"
)

func main() {
	blockchain.BlockChain()
	defer db.Close()
}
