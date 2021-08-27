package blockchain

import (
	"fmt"
	"strings"
	"time"

	"github.com/chiwon99881/one/db"
	"github.com/chiwon99881/one/utils"
)

type Block struct {
	Transactions []*Tx  `json:"transactions"`
	Hash         string `json:"hash"`
	PrevHash     string `json:"prevHash,omitempty"`
	Height       int    `json:"height"`
	Timestamp    int    `json:"-"`
	Nouce        int    `json:"nounce"`
	Difficulty   int    `json:"difficulty"`
}

const (
	currentDifficulty int = 2
)

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

func FindBlock(hash string) *Block {
	block := &Block{}
	blockAsBytes := db.FindBlock(hash)
	utils.FromBytes(block, blockAsBytes)
	return block
}

func persistBlock(newBlock *Block) {
	blockAsBytes := utils.ToBytes(newBlock)
	db.SaveBlockDB(newBlock.Hash, blockAsBytes)
}

func (b *Block) mine() {
	currentPreFix := strings.Repeat("0", currentDifficulty)
	// mine reward
	for {
		hashAsBytes := utils.ToBytes(b)
		hash := fmt.Sprintf("%x", hashAsBytes)
		done := strings.HasPrefix(hash, currentPreFix)
		if done {
			b.Hash = hash
			b.Timestamp = int(time.Now().Unix())
			break
		}
		b.Nouce++
	}
}

// CreateBlock is generate new block.
func CreateBlock(prevHash string, height int) *Block {
	b := &Block{
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Nouce:      0,
		Difficulty: currentDifficulty,
	}
	b.Transactions = m.TxToConfirm()
	b.mine()
	return b
}
