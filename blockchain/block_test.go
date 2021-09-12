package blockchain

import (
	"sync"
	"testing"

	"github.com/chiwon99881/one/utils"
)

func TestBlocks(t *testing.T) {
	t.Run("find block is nil, blocks length should be zero", func(t *testing.T) {
		chain = &blockchain{
			NewestHash: "",
		}

		dbOperator = &testDataBase{
			testFindBlock: func(hash string) []byte {
				return nil
			},
		}
		blocks := Blocks(chain)
		if len(blocks) != 0 {
			t.Fatalf("blocks length should be 0 but got %d", len(blocks))
		}
		chain = nil
	})
	t.Run("find block is not nil, blocks length should not be zero", func(t *testing.T) {
		index := 0
		chain = &blockchain{
			NewestHash: "",
		}
		dbOperator = &testDataBase{
			testFindBlock: func(hash string) []byte {
				defer func() {
					index++
				}()
				var block *Block
				if index == 0 {
					block = &Block{
						Hash:     "first",
						PrevHash: "hash",
					}
				} else {
					block = &Block{
						Hash:     "hash",
						PrevHash: "",
					}
				}
				blockAsBytes := utils.ToBytes(block)
				return blockAsBytes
			},
		}
		blocks := Blocks(chain)
		if len(blocks) != 2 {
			t.Fatalf("blocks length should be 2 but got %d", len(blocks))
		}
		chain = nil
	})
}

func TestFindBlock(t *testing.T) {
	dbOperator = &testDataBase{
		testFindBlock: func(hash string) []byte {
			block := &Block{
				Hash:   hash,
				Height: 1,
			}
			blockBytes := utils.ToBytes(block)
			return blockBytes
		},
	}
	block := FindBlock("hash")
	if block.Hash != "hash" {
		t.Fatalf("block hash should be 'hash' but got %s", block.Hash)
	}
	if block.Height != 1 {
		t.Fatalf("block height should be 1 but got %d", block.Height)
	}
}

func TestPersistBlock(t *testing.T) {
	newBlock := &Block{
		Hash: "newHash",
	}
	saveBlockResult := &Block{}
	dbOperator = &testDataBase{
		testSaveBlockDB: func(key string, data []byte) string {
			utils.FromBytes(saveBlockResult, data)
			return saveBlockResult.Hash
		},
	}
	persistBlock(newBlock)
	if saveBlockResult.Hash != "newHash" {
		t.Fatalf("persist block's hash should be 'newHash' but got %s", saveBlockResult.Hash)
	}
}

func TestMine(t *testing.T) {
	block := &Block{
		Difficulty: 2,
		Nounce:     0,
		Hash:       "",
	}
	block.mine()

	if block.Hash == "" {
		t.Fatalf("block hash should be change from initialized hash")
	}
	if block.Nounce == 0 {
		t.Fatalf("block's Nounce should be change from initialized nounce")
	}
}

func TestCreateBlock(t *testing.T) {
	chain = nil
	once = *new(sync.Once)
	dbOperator = &testDataBase{
		testGetExistChain: func() []byte {
			testChain := &blockchain{
				NewestHash:        "",
				CurrentDifficulty: 2,
			}
			chainBytes := utils.ToBytes(testChain)
			return chainBytes
		},
		testFindBlock: func(hash string) []byte {
			return nil
		},
	}
	m = &mempool{
		Txs: make(map[string]*Tx),
	}
	tx := &Tx{}
	m.Txs["hash"] = tx
	block := CreateBlock("prevHash", 6)
	if block.Height != 6 {
		t.Fatalf("block's height should be 6 but got %d", block.Height)
	}
	if block.Difficulty != 2 {
		t.Fatalf("block's difficulty should be 2 but got %d", block.Difficulty)
	}
	if len(block.Transactions) != 2 {
		t.Fatalf("block's txs length should be 2 but got %d", len(block.Transactions))
	}
	m = nil
}
