package blockchain

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/chiwon99881/one/db"
	"github.com/chiwon99881/one/utils"
)

type Block struct {
	Data      string `json:"data"`
	Hash      string `json:"hash"`
	PrevHash  string `json:"prevHash,omitempty"`
	Height    int    `json:"height"`
	Timestamp int    `json:"-"`
}

func Blocks(bc *chain) []*Block {
	var blocks []*Block
	hashCursor := bc.NewestHash
	for {
		block := &Block{}
		blockAsBytes := db.FindBlock(hashCursor)
		if blockAsBytes == nil {
			break
		}
		utils.FromBytes(block, blockAsBytes)
		blocks = append(blocks, block)
		if block.PrevHash == "" {
			break
		}
		hashCursor = block.PrevHash
	}
	return blocks
}

func persistBlock(newBlock *Block) {
	blockAsBytes := utils.ToBytes(newBlock)
	db.SaveBlockDB(newBlock.Hash, blockAsBytes)
}

func (b *Block) hash() {
	blockAsBytes := utils.ToBytes(b)
	bytes := sha256.Sum256(blockAsBytes)
	hash := fmt.Sprintf("%x", bytes)
	b.Hash = hash
}

// CreateBlock is generate new block.
func CreateBlock(data, prevHash string, height int) *Block {
	b := &Block{
		Hash:      "",
		Data:      data,
		PrevHash:  prevHash,
		Height:    height,
		Timestamp: int(time.Now().Unix()),
	}
	b.hash()
	return b
}
