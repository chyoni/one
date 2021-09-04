package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"

	"github.com/chiwon99881/one/db"
	"github.com/chiwon99881/one/utils"
	"github.com/chiwon99881/one/wallet"
)

type Tx struct {
	TxID   string   `json:"txID"`
	TxIns  []*TxIn  `json:"txIns"`
	TxOuts []*TxOut `json:"txOuts"`
}

type TxIn struct {
	TxID      string `json:"txID"`
	Index     int    `json:"index"`
	Amount    int    `json:"amount"`
	Signature string `json:"signature"`
}

type TxOut struct {
	Amount  int    `json:"amount"`
	Address string `json:"address"`
}

type UTxOut struct {
	TxID    string `json:"txID"`
	Index   int    `json:"index"`
	Amount  int    `json:"amount"`
	Address string `json:"address"`
}

type mempool struct {
	Txs []*Tx
	m   sync.Mutex
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

func FindTx(txID string) *Tx {
	txs := Txs()
	for _, tx := range txs {
		if tx.TxID == txID {
			return tx
		}
	}
	return nil
}

func GetBalanceByAddress(address string) int {
	total := 0
	for _, uTxOut := range GetUTxOutsByAddress(address) {
		total += uTxOut.Amount
	}
	return total
}

func GetUTxOutsByAddress(address string) []*UTxOut {
	var ownedUTxOuts []*UTxOut
	sTxOut := make(map[string]bool)
	txs := Txs()
	txs = append(txs, Mempool().Txs...)
	for _, tx := range txs {
		for _, txIn := range tx.TxIns {
			if txIn.Signature == "COINBASE" {
				break
			}
			if FindTx(txIn.TxID).TxOuts[txIn.Index].Address == address {
				sTxOut[txIn.TxID] = true
			}
		}
		for index, txOut := range tx.TxOuts {
			if txOut.Address == address {
				_, isTrue := sTxOut[tx.TxID]
				if !isTrue && !isOnMempool(tx.TxID, address) {
					uTxOut := &UTxOut{
						TxID:    tx.TxID,
						Index:   index,
						Amount:  txOut.Amount,
						Address: txOut.Address,
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

func (tx *Tx) verifyTx() bool {
	for _, txIn := range tx.TxIns {
		tx := FindTx(txIn.TxID)
		if tx == nil {
			return false
		}
		publicKey := tx.TxOuts[txIn.Index].Address
		verified := wallet.Verify(publicKey, txIn.Signature, txIn.TxID)
		if !verified {
			return false
		}
	}
	return true
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
			Amount:    uTxOut.Amount,
			Index:     uTxOut.Index,
			TxID:      uTxOut.TxID,
			Signature: wallet.Sign(uTxOut.TxID),
		}
		txIns = append(txIns, shiftTxIn)
		total += uTxOut.Amount
	}
	toTxOut := &TxOut{
		Address: to,
		Amount:  amount,
	}
	txOuts = append(txOuts, toTxOut)
	if total >= amount {
		exchange := total - amount
		exchangeTxOut := &TxOut{
			Address: from,
			Amount:  exchange,
		}
		txOuts = append(txOuts, exchangeTxOut)
	}
	tx := &Tx{
		TxIns:  txIns,
		TxOuts: txOuts,
	}
	tx.generateTxID()
	result := tx.verifyTx()
	if !result {
		return nil, errors.New("this signature invalid")
	}
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) (*Tx, error) {
	tx, err := makeTx(wallet.Wallet().Address, to, amount)
	if err != nil {
		return nil, err
	}
	var newestTxs []*Tx
	newestTxs = append(newestTxs, tx)
	newestTxs = append(newestTxs, m.Txs...)
	m.Txs = newestTxs
	mBytes := utils.ToBytes(m)
	db.PushOnMempool(mBytes)
	return tx, nil
}

func (m *mempool) TxToConfirm() []*Tx {
	var txs []*Tx
	txs = append(txs, m.Txs...)
	m.Txs = nil
	mBytes := utils.ToBytes(m)
	db.PushOnMempool(mBytes)
	return txs
}

func isOnMempool(txID, address string) bool {
	isOn := false
Outer:
	for _, tx := range Mempool().Txs {
		for _, txIn := range tx.TxIns {
			if txIn.TxID == txID && FindTx(txIn.TxID).TxOuts[txIn.Index].Address == address {
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
		TxID:      "COINBASE",
		Index:     -1,
		Amount:    50,
		Signature: "COINBASE",
	}
	txOut := &TxOut{
		Address: wallet.Wallet().Address,
		Amount:  50,
	}
	txIns = append(txIns, txIn)
	txOuts = append(txOuts, txOut)
	coinbaseTx := &Tx{
		TxIns:  txIns,
		TxOuts: txOuts,
	}
	coinbaseTxAsBytes := utils.ToBytes(coinbaseTx)
	coinbaseTxAsBytes = append(coinbaseTxAsBytes, utils.ToBytes(chain.Height)...)
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
