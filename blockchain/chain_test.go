package blockchain

import "testing"

type testDataBase struct {
	testGetExistChain func() []byte
}

func (t *testDataBase) GetExistChain() []byte {
	return t.testGetExistChain()
}
func TestBlockChain(t *testing.T) {
	dbOperator = &testDataBase{
		testGetExistChain: func() []byte {
			return nil
		},
	}
	chain := BlockChain()
	if chain.Height != 0 || chain.NewestHash != "" || chain.CurrentDifficulty != DefaultDifficulty {
		t.Fatalf("chain's height should be 0, but got %d", chain.Height)
	}
}
