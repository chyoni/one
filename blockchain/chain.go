package blockchain

import (
	"sync"

	"github.com/chiwon99881/one/db"
	"github.com/chiwon99881/one/utils"
)

type chain struct {
	Height            int    `json:"height"`
	NewestHash        string `json:"newestHash"`
	CurrentDifficulty int    `json:"currentDifficulty"`
}

var once sync.Once
var blockchain *chain

const (
	DefaultDifficulty   int = 2
	ReCalculateInterval int = 5
	IncreaseDifficulty  int = 5
	DecreaseDifficulty  int = 10
)

func (bc *chain) persistChain(newBlock *Block) {
	bc.Height++
	bc.NewestHash = newBlock.Hash
	chainAsBytes := utils.ToBytes(bc)
	db.SaveChainDB(chainAsBytes)
}

func restoreChain(data []byte) {
	utils.FromBytes(blockchain, data)
}

func AddBlock(bc *chain) {
	block := CreateBlock(bc.NewestHash, bc.Height+1)
	persistBlock(block)
	bc.persistChain(block)
}

func GetCurrentDifficulty(blockchain *chain) int {
	if len(Blocks(blockchain)) == 0 {
		return blockchain.CurrentDifficulty
	}
	if len(Blocks(blockchain))%5 == 0 {
		return blockchain.reCalculateDifficulty()
	} else {
		return blockchain.CurrentDifficulty
	}
}

func (blockchain *chain) reCalculateDifficulty() int {
	allBlocks := Blocks(blockchain)
	latestBlock := allBlocks[0]
	latestRecalculateBlock := allBlocks[ReCalculateInterval-1]
	interval := (latestBlock.Timestamp / 60) - (latestRecalculateBlock.Timestamp / 60)
	if interval < IncreaseDifficulty {
		blockchain.CurrentDifficulty += 1
		return blockchain.CurrentDifficulty
	} else if interval > DecreaseDifficulty {
		blockchain.CurrentDifficulty -= 1
		return blockchain.CurrentDifficulty
	} else {
		return blockchain.CurrentDifficulty
	}
}

func BlockChain() *chain {
	if blockchain == nil {
		once.Do(func() {
			blockchain = &chain{
				Height:            0,
				NewestHash:        "",
				CurrentDifficulty: DefaultDifficulty,
			}
			existChain := db.GetExistChain()
			if existChain != nil {
				restoreChain(existChain)
			}
		})
	}
	return blockchain
}
