package blockchain

import (
	"sync"

	"github.com/chiwon99881/one/db"
	"github.com/chiwon99881/one/utils"
)

type chain struct {
	Height     int
	NewestHash string
}

var once sync.Once
var blockchain *chain

func (bc *chain) persistChain(newBlock *Block) {
	bc.Height++
	bc.NewestHash = newBlock.Hash
	hashAsBytes := utils.ToBytes(bc.NewestHash)
	db.SaveChainDB(hashAsBytes)
}

func AddBlock(bc *chain, data string) {
	block := CreateBlock(data, bc.NewestHash, bc.Height+1)
	persistBlock(block)
	bc.persistChain(block)
}

func BlockChain() *chain {
	if blockchain == nil {
		once.Do(func() {
			blockchain = &chain{
				Height:     0,
				NewestHash: "",
			}
			AddBlock(blockchain, "first block")
		})
	}
	return blockchain
}
