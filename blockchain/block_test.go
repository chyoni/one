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
			testFindBlock: func() []byte {
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
			testFindBlock: func() []byte {
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
