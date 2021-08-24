package blockchain

import (
	"errors"
	"sync"
)

type Tx struct {
	TxIns  []*TxIn  `json:"txIns"`
	TxOuts []*TxOut `json:"txOuts"`
}

type TxIn struct {
	Amount int    `json:"amount"`
	Owner  string `json:"owner"`
}

type TxOut struct {
	Amount int    `json:"amount"`
	Owner  string `json:"address"`
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
	for _, tx := range Txs() {
		for _, txOut := range tx.TxOuts {
			if txOut.Owner == from {
				total += txOut.Amount
			}
		}
	}
	return total
}

func getTxOutByAddress(from string) []*TxOut {
	var ownedTxOuts []*TxOut
	txs := Txs()
	for _, tx := range txs {
		for _, txOut := range tx.TxOuts {
			if txOut.Owner == from {
				ownedTxOuts = append(ownedTxOuts, txOut)
			}
		}
	}
	return ownedTxOuts
}

func makeTx(from, to string, amount int) (*Tx, error) {
	if GetBalanceByAddress(from) < amount {
		return nil, errors.New("not enough money")
	}
	var txIns []*TxIn
	var txOuts []*TxOut
	var total int
	for _, txOut := range getTxOutByAddress(from) {
		if total >= amount {
			break
		}
		shiftTxIn := &TxIn{
			Owner:  txOut.Owner,
			Amount: txOut.Amount,
		}
		txIns = append(txIns, shiftTxIn)
		total += txOut.Amount
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
	return tx, nil
}

func (m *mempool) AddTx(from, to string, amount int) error {
	tx, err := makeTx(from, to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func Mempool() *mempool {
	if m == nil {
		memOnce.Do(func() {
			m = &mempool{}
		})
	}
	return m
}
