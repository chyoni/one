package blockchain

import (
	"fmt"
	"sync"
	"testing"

	"github.com/chiwon99881/one/utils"
)

type testDataBase struct {
	testGetExistChain func() []byte
	testFindBlock     func(hash string) []byte
}

func (t *testDataBase) GetExistChain() []byte {
	return t.testGetExistChain()
}

func (testDataBase) SaveChainDB(data []byte) {
	fmt.Println("executing saveChainDB func for test")
}

func (t *testDataBase) FindBlock(hash string) []byte {
	return t.testFindBlock(hash)
}

func TestBlockChain(t *testing.T) {
	t.Run("get exist chain is nil", func(t *testing.T) {
		dbOperator = &testDataBase{
			testGetExistChain: func() []byte {
				return nil
			},
		}
		chain := BlockChain()
		if chain.Height != 0 {
			t.Fatalf("chain's height should be 0, but got %d", chain.Height)
		}
		if chain.NewestHash != "" {
			t.Fatalf("chain's NewestHash should be '' but got %s", chain.NewestHash)
		}
		if chain.CurrentDifficulty != DefaultDifficulty {
			t.Fatalf("chain's CurrentDifficulty is 2 but got %d", chain.CurrentDifficulty)
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

func TestPersistChain(t *testing.T) {
	testBlock := &Block{
		Hash: "",
	}
	testChain := &blockchain{
		Height:     1,
		NewestHash: "abcd",
	}

	dbOperator = &testDataBase{}

	testChain.persistChain(testBlock)
	if testChain.Height != 2 {
		t.Fatalf("chain's height should be 2 but got %d", testBlock.Height)
	}
	if testChain.NewestHash != "" {
		t.Fatalf("chain's NewestHash should be '' but got %s", testChain.NewestHash)
	}
}
