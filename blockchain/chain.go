package blockchain

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/chiwon99881/one/db"
	"github.com/chiwon99881/one/interfaces"
	"github.com/chiwon99881/one/utils"
)

type blockchain struct {
	Height            int    `json:"height"`
	NewestHash        string `json:"newestHash"`
	CurrentDifficulty int    `json:"currentDifficulty"`
	m                 sync.Mutex
}

var dbOperator interfaces.DBOperation = db.DataBase{}

var once sync.Once
var chain *blockchain

const (
	DefaultDifficulty   int = 2
	ReCalculateInterval int = 5
	IncreaseDifficulty  int = 5
	DecreaseDifficulty  int = 10
)

func (chain *blockchain) persistChain(newBlock *Block) {
	chain.Height++
	chain.NewestHash = newBlock.Hash
	chainAsBytes := utils.ToBytes(chain)
	db.SaveChainDB(chainAsBytes)
}

func restoreChain(data []byte) {
	utils.FromBytes(chain, data)
}

func AddBlock(chain *blockchain) *Block {
	block := CreateBlock(chain.NewestHash, chain.Height+1)
	persistBlock(block)
	chain.persistChain(block)
	return block
}

func GetCurrentDifficulty(chain *blockchain) int {
	if len(Blocks(chain)) == 0 {
		return chain.CurrentDifficulty
	}
	if len(Blocks(chain))%5 == 0 {
		return chain.reCalculateDifficulty()
	} else {
		return chain.CurrentDifficulty
	}
}

func (chain *blockchain) reCalculateDifficulty() int {
	allBlocks := Blocks(chain)
	latestBlock := allBlocks[0]
	latestRecalculateBlock := allBlocks[ReCalculateInterval-1]
	interval := (latestBlock.Timestamp / 60) - (latestRecalculateBlock.Timestamp / 60)
	if interval < IncreaseDifficulty {
		chain.CurrentDifficulty += 1
		return chain.CurrentDifficulty
	} else if interval > DecreaseDifficulty {
		chain.CurrentDifficulty -= 1
		return chain.CurrentDifficulty
	} else {
		return chain.CurrentDifficulty
	}
}

func Status(rw http.ResponseWriter) error {
	chain.m.Lock()
	defer chain.m.Unlock()
	err := json.NewEncoder(rw).Encode(chain)
	return err
}

func GetNewestHash() string {
	return chain.NewestHash
}

func BlockChain() *blockchain {
	if chain == nil {
		once.Do(func() {
			chain = &blockchain{
				Height:            0,
				NewestHash:        "",
				CurrentDifficulty: DefaultDifficulty,
			}
			existChain := dbOperator.GetExistChain()
			if existChain != nil {
				restoreChain(existChain)
			}
		})
	}
	return chain
}
