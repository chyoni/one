package blockchain

import (
	"sync"
	"testing"

	"github.com/chiwon99881/one/utils"
)

type testDataBase struct {
	testGetExistChain func() []byte
}

func (t *testDataBase) GetExistChain() []byte {
	return t.testGetExistChain()
}
func TestBlockChain(t *testing.T) {
	t.Run("get exist chain is nil", func(t *testing.T) {
		dbOperator = &testDataBase{
			testGetExistChain: func() []byte {
				return nil
			},
		}
		chain := BlockChain()
		if chain.Height != 0 || chain.NewestHash != "" || chain.CurrentDifficulty != DefaultDifficulty {
			t.Fatalf("chain's height should be 0, but got %d", chain.Height)
		}
	})

	t.Run("get exist chain is not nil", func(t *testing.T) {
		chain = nil
		once = *new(sync.Once)
		dbOperator = &testDataBase{
			testGetExistChain: func() []byte {
				testChain := &blockchain{
					Height:            10,
					NewestHash:        "",
					CurrentDifficulty: 5,
				}
				chainBytes := utils.ToBytes(testChain)
				return chainBytes
			},
		}
		tChain := BlockChain()
		if tChain.Height != 10 || tChain.NewestHash != "" || tChain.CurrentDifficulty != 5 {
			t.Fatalf("chain's height should be 10 but got %d", tChain.Height)
		}
	})
}
