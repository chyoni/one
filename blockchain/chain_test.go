package blockchain

import (
	"fmt"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"

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

func TestGetCurrentDifficulty(t *testing.T) {
	t.Run("get current chain difficulty", func(t *testing.T) {
		index := 0
		chain = &blockchain{
			NewestHash:        "firstHash",
			CurrentDifficulty: 2,
		}
		dbOperator = &testDataBase{
			testFindBlock: func(hash string) []byte {
				defer func() {
					index++
				}()
				var block *Block
				if index == 0 {
					block = &Block{
						Hash:     hash,
						PrevHash: "prev",
					}
				} else {
					block = &Block{
						Hash:     "prev",
						PrevHash: "",
					}
				}
				blockAsBytes := utils.ToBytes(block)
				return blockAsBytes
			},
		}
		diff := GetCurrentDifficulty(chain)
		if diff != chain.CurrentDifficulty {
			t.Fatalf("current difficulty should be %d but got %d", chain.CurrentDifficulty, diff)
		}
	})
}

func TestReCalculateDifficulty(t *testing.T) {
	t.Run("get increased difficulty", func(t *testing.T) {
		index := 0
		chain = &blockchain{
			NewestHash:        "hash",
			CurrentDifficulty: 2,
		}
		dbOperator = &testDataBase{
			testFindBlock: func(hash string) []byte {
				if index == 5 {
					index = 0
					return nil
				}
				block := &Block{
					Hash:      "hash1",
					PrevHash:  "prevHash1",
					Timestamp: int(time.Now().Unix()) + (4 - index),
				}
				blockAsBytes := utils.ToBytes(block)
				index++
				return blockAsBytes
			},
		}
		GetCurrentDifficulty(chain)
		if chain.CurrentDifficulty != 3 {
			t.Fatalf("chain's CurrentDifficulty should be 3 but got %d", chain.CurrentDifficulty)
		}
	})

	t.Run("get decreased difficulty", func(t *testing.T) {
		index := 0
		chain = &blockchain{
			NewestHash:        "hash",
			CurrentDifficulty: 2,
		}
		dbOperator = &testDataBase{
			testFindBlock: func(hash string) []byte {
				if index == 5 {
					index = 0
					return nil
				}
				var block *Block
				switch index {
				case 0:
					block = &Block{
						Hash:      "hash1",
						PrevHash:  "prevHash1",
						Timestamp: int(time.Now().Unix()) + (1000000),
					}
				default:
					block = &Block{
						Hash:      "hash1",
						PrevHash:  "prevHash1",
						Timestamp: int(time.Now().Unix()),
					}
				}
				blockAsBytes := utils.ToBytes(block)
				index++
				return blockAsBytes
			},
		}
		GetCurrentDifficulty(chain)
		if chain.CurrentDifficulty != 1 {
			t.Fatalf("chain's CurrentDifficulty should be 1 but got %d", chain.CurrentDifficulty)
		}
	})

	t.Run("get same difficulty", func(t *testing.T) {
		index := 0
		chain = &blockchain{
			NewestHash:        "hash",
			CurrentDifficulty: 2,
		}
		dbOperator = &testDataBase{
			testFindBlock: func(hash string) []byte {
				if index == 5 {
					index = 0
					return nil
				}
				var block *Block
				switch index {
				case 0:
					block = &Block{
						Hash:      "hash1",
						PrevHash:  "prevHash1",
						Timestamp: int(time.Now().Unix()) + (420),
					}
				case 4:
					block = &Block{
						Hash:      "hash4",
						PrevHash:  "prevHash4",
						Timestamp: int(time.Now().Unix()),
					}
				default:
					block = &Block{
						Hash:      "hash" + strconv.Itoa(index),
						PrevHash:  "prevHash" + strconv.Itoa(index),
						Timestamp: int(time.Now().Unix()),
					}
				}
				blockAsBytes := utils.ToBytes(block)
				index++
				return blockAsBytes
			},
		}
		GetCurrentDifficulty(chain)
		if chain.CurrentDifficulty != 2 {
			t.Fatalf("chain's CurrentDifficulty should be 2 but got %d", chain.CurrentDifficulty)
		}
	})
}

func TestStatus(t *testing.T) {
	chain = &blockchain{
		NewestHash:        "hash",
		Height:            5,
		CurrentDifficulty: 3,
	}
	rw := httptest.NewRecorder()
	err := Status(rw)
	if err != nil {
		t.Fatalf("chain status should be OK but got %s", err)
	}
}
