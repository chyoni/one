package blockchain

import (
	"github.com/chiwon99881/one/db"
	"github.com/chiwon99881/one/utils"
)

func HandleSendAllBlocksMessage(blocks []*Block) {
	newestBlock := blocks[0]
	chain.m.Lock()
	defer chain.m.Unlock()

	chain.CurrentDifficulty = newestBlock.Difficulty
	chain.Height = newestBlock.Height
	chain.NewestHash = newestBlock.Hash

	chainAsBytes := utils.ToBytes(chain)

	db.CreateAfterDeleteDB()
	db.SaveChainDB(chainAsBytes)
	for _, block := range blocks {
		blockAsBytes := utils.ToBytes(block)
		db.SaveBlockDB(block.Hash, blockAsBytes)
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
}

func HandleNewTxMessage(tx *Tx) {
	Mempool().m.Lock()
	defer Mempool().m.Unlock()

	tempTxs := make(map[string]*Tx)
	tempTxs[tx.TxID] = tx
	for key, value := range Mempool().Txs {
		tempTxs[key] = value
	}
	m.Txs = tempTxs
	mBytes := utils.ToBytes(m)
	db.PushOnMempool(mBytes)
}
