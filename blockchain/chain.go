package blockchain

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/chiwon99881/one/db"
	"github.com/chiwon99881/one/utils"
)

type blockchain struct {
	Height            int    `json:"height"`
	NewestHash        string `json:"newestHash"`
	CurrentDifficulty int    `json:"currentDifficulty"`
	m                 sync.Mutex
}

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

func AddBlock(chain *blockchain) {
	block := CreateBlock(chain.NewestHash, chain.Height+1)
	persistBlock(block)
	chain.persistChain(block)
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

func HandleSendAllBlocksMessage(blocks []*Block) {
	newestBlock := blocks[0]
	chain.m.Lock()
	defer chain.m.Unlock()

	chain.CurrentDifficulty = newestBlock.Difficulty
	chain.Height = newestBlock.Height
	chain.NewestHash = newestBlock.Hash

	chainAsBytes := utils.ToBytes(chain)

	db.CreateAfterDeleteDB()
	db.SaveChainDB(chainAsBytes)
	for _, block := range blocks {
		blockAsBytes := utils.ToBytes(block)
		db.SaveBlockDB(block.Hash, blockAsBytes)
	}
}

func Status(rw http.ResponseWriter) error {
	chain.m.Lock()
	defer chain.m.Unlock()
	err := json.NewEncoder(rw).Encode(chain)
	return err
}

func BlockChain() *blockchain {
	if chain == nil {
		once.Do(func() {
			chain = &blockchain{
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
	return chain
}
