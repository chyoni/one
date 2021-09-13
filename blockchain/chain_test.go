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
	testSaveBlockDB   func(key string, data []byte) string
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

func (t *testDataBase) SaveBlockDB(key string, data []byte) {
	t.testSaveBlockDB(key, data)
}

func (testDataBase) PushOnMempool(data []byte) {
	fmt.Println("executing pushOnMempool func for test")
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

func TestAddBlock(t *testing.T) {
	saveBlockResult := &Block{}
	dbOperator = &testDataBase{
		testFindBlock: func(hash string) []byte {
			return nil
		},
		testSaveBlockDB: func(key string, data []byte) string {
			utils.FromBytes(saveBlockResult, data)
			return saveBlockResult.Hash
		},
	}
	m = &mempool{
		Txs: make(map[string]*Tx),
	}
	tx := &Tx{}
	m.Txs["hash"] = tx

	chain := &blockchain{
		NewestHash: "hash",
		Height:     1,
	}
	block := AddBlock(chain)

	if block.Hash != saveBlockResult.Hash {
		t.Fatalf("block hash is should be same saveBlockResult.Hash")
	}
	if chain.Height != 2 {
		t.Fatalf("chain's height should be 2 but got %d", chain.Height)
	}
	if chain.NewestHash != block.Hash {
		t.Fatalf("chain's newestHash should be same block.Hash")
	}
	m = nil
}
