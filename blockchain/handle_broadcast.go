package blockchain

import (
	"github.com/chiwon99881/one/db"
	"github.com/chiwon99881/one/utils"
)

const (
	ChangeMyBlockChain int = iota
	SendMyBlockChain
	NothingToDoAnything
)

func HandleSendNewestBlockMessage(block *Block) int {
	BlockChain().m.Lock()
	defer BlockChain().m.Unlock()

	myChain := BlockChain()
	if block.Height > myChain.Height {
		return ChangeMyBlockChain
	} else if block.Height < myChain.Height {
		return SendMyBlockChain
	} else {
		return NothingToDoAnything
	}
}

func HandleSendAllBlocksMessage(blocks []*Block) {
	newestBlock := blocks[0]
	BlockChain().m.Lock()
	defer BlockChain().m.Unlock()

	chain.CurrentDifficulty = newestBlock.Difficulty
	chain.Height = newestBlock.Height
	chain.NewestHash = newestBlock.Hash

	chainAsBytes := utils.ToBytes(chain)

	db.CreateAfterDeleteDB()
	dbOperator.SaveChainDB(chainAsBytes)
	for _, block := range blocks {
		blockAsBytes := utils.ToBytes(block)
		dbOperator.SaveBlockDB(block.Hash, blockAsBytes)
	}
}

func HandleNewBlockMessage(block *Block) {
	chain.m.Lock()
	defer chain.m.Unlock()

	persistBlock(block)
	chain.persistChain(block)

	for _, tx := range block.Transactions {
		for key := range Mempool().Txs {
			if tx.TxID == key {
				delete(Mempool().Txs, key)
			}
		}
	}
	mBytes := utils.ToBytes(Mempool())
	db.PushOnMempool(mBytes)
}

func HandleNewTxMessage(tx *Tx) {
	Mempool().m.Lock()
	defer Mempool().m.Unlock()

	tempTxs := make(map[string]*Tx)
	tempTxs[tx.TxID] = tx
	for key, value := range Mempool().Txs {
		tempTxs[key] = value
	}
	Mempool().Txs = tempTxs
	mBytes := utils.ToBytes(Mempool())
	db.PushOnMempool(mBytes)
}
