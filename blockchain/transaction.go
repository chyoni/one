package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"

	"github.com/chiwon99881/one/db"
	"github.com/chiwon99881/one/utils"
)

type Tx struct {
	TxID   string   `json:"txID"`
	TxIns  []*TxIn  `json:"txIns"`
	TxOuts []*TxOut `json:"txOuts"`
}

type TxIn struct {
	TxID   string `json:"txID"`
	Index  int    `json:"index"`
	Amount int    `json:"amount"`
	Owner  string `json:"owner"`
}

type TxOut struct {
	Amount int    `json:"amount"`
	Owner  string `json:"owner"`
}

type UTxOut struct {
	TxID   string `json:"txID"`
	Index  int    `json:"index"`
	Amount int    `json:"amount"`
	Owner  string `json:"owner"`
}

type mempool struct {
	Txs []*Tx
}

var m *mempool
var memOnce sync.Once

func Txs() []*Tx {
	var txs []*Tx
	blocks := Blocks(BlockChain())
	for _, block := range blocks {
		txs = append(txs, block.Transactions...)
	}
	return txs
}

func GetBalanceByAddress(from string) int {
	total := 0
	for _, uTxOut := range GetUTxOutsByAddress(from) {
		total += uTxOut.Amount
	}
	return total
}

func GetUTxOutsByAddress(from string) []*UTxOut {
	var ownedUTxOuts []*UTxOut
	sTxOut := make(map[string]bool)
	txs := Txs()
	for _, tx := range txs {
		for _, txIn := range tx.TxIns {
			if txIn.Owner == from {
				sTxOut[txIn.TxID] = true
			}
		}
		for index, txOut := range tx.TxOuts {
			if txOut.Owner == from {
				_, isTrue := sTxOut[tx.TxID]
				if !isTrue && !isOnMempool(tx.TxID, from) {
					uTxOut := &UTxOut{
						TxID:   tx.TxID,
						Index:  index,
						Amount: txOut.Amount,
						Owner:  txOut.Owner,
					}
					ownedUTxOuts = append(ownedUTxOuts, uTxOut)
				}
			}
		}
	}
	return ownedUTxOuts
}

func (tx *Tx) generateTxID() {
	IDAsBytes := utils.ToBytes(tx)
	txID := fmt.Sprintf("%x", sha256.Sum256(IDAsBytes))
	tx.TxID = txID
}

func makeTx(from, to string, amount int) (*Tx, error) {
	if GetBalanceByAddress(from) < amount {
		return nil, errors.New("not enough money")
	}
	var txIns []*TxIn
	var txOuts []*TxOut
	var total int
	for _, uTxOut := range GetUTxOutsByAddress(from) {
		if total >= amount {
			break
		}
		shiftTxIn := &TxIn{
			Owner:  uTxOut.Owner,
			Amount: uTxOut.Amount,
			Index:  uTxOut.Index,
			TxID:   uTxOut.TxID,
		}
		txIns = append(txIns, shiftTxIn)
		total += uTxOut.Amount
	}
	toTxOut := &TxOut{
		Owner:  to,
		Amount: amount,
	}
	txOuts = append(txOuts, toTxOut)
	if total >= amount {
		exchange := total - amount
		exchangeTxOut := &TxOut{
			Owner:  from,
			Amount: exchange,
		}
		txOuts = append(txOuts, exchangeTxOut)
	}
	tx := &Tx{
		TxIns:  txIns,
		TxOuts: txOuts,
	}
	tx.generateTxID()
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("chyonee", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	mBytes := utils.ToBytes(m)
	db.PushOnMempool(mBytes)
	return nil
}

func (m *mempool) TxToConfirm() []*Tx {
	var txs []*Tx
	txs = append(txs, m.Txs...)
	m.Txs = nil
	mBytes := utils.ToBytes(m)
	db.PushOnMempool(mBytes)
	return txs
}

func isOnMempool(txID, owner string) bool {
	isOn := false
Outer:
	for _, tx := range m.Txs {
		for _, txIn := range tx.TxIns {
			if txIn.TxID == txID && txIn.Owner == owner {
				isOn = true
				break Outer
			}
		}
	}
	return isOn
}

func (b *Block) coinbaseTx() {
	var txIns []*TxIn
	var txOuts []*TxOut
	txIn := &TxIn{
		TxID:   "COINBASE",
		Index:  -1,
		Owner:  "COINBASE",
		Amount: 50,
	}
	txOut := &TxOut{
		Owner:  "chyonee",
		Amount: 50,
	}
	txIns = append(txIns, txIn)
	txOuts = append(txOuts, txOut)
	coinbaseTx := &Tx{
		TxIns:  txIns,
		TxOuts: txOuts,
	}
	coinbaseTxAsBytes := utils.ToBytes(coinbaseTx)
	coinbaseTxAsBytes = append(coinbaseTxAsBytes, utils.ToBytes(blockchain.Height)...)
	coinbaseTx.TxID = fmt.Sprintf("%x", sha256.Sum256(coinbaseTxAsBytes))
	b.Transactions = append(b.Transactions, coinbaseTx)
}

func Mempool() *mempool {
	if m == nil {
		memOnce.Do(func() {
			m = &mempool{}
			memData := db.GetExistMempool()
			if memData != nil {
				utils.FromBytes(m, memData)
			}
		})
	}
	return m
}
