package blockchain

import (
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
