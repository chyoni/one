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
}

func HandleNewTxMessage(tx *Tx) {
	Mempool().m.Lock()
	defer Mempool().m.Unlock()

	var newestTxs []*Tx
	newestTxs = append(newestTxs, tx)
	newestTxs = append(newestTxs, Mempool().Txs...)
	m.Txs = newestTxs
	mBytes := utils.ToBytes(m)
	db.PushOnMempool(mBytes)
}
