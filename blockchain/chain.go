package blockchain

import "sync"

type chain struct {
	Blocks     []*Block
	NewestHash string
}

var once sync.Once
var blockchain *chain

func AddBlock(bc *chain) {
	block := CreateBlock("first Block", bc.NewestHash)
	bc.Blocks = append(bc.Blocks, block)
	bc.NewestHash = block.Hash
}

func BlockChain() *chain {
	if blockchain == nil {
		once.Do(func() {
			blockchain = &chain{}
			AddBlock(blockchain)
		})
	}
	return blockchain
}
