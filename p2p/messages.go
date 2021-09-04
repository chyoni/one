package p2p

import (
	"fmt"

	"github.com/chiwon99881/one/blockchain"
	"github.com/chiwon99881/one/utils"
)

type MessageKind int

type Message struct {
	MessageKind MessageKind `json:"messageKind"`
	Payload     []byte      `json:"payload"`
}

const (
	SendNewestBlockMessage MessageKind = iota
	RequestAllBlocksMessage
	SendAllBlocksMessage
	NewTransactionMessage
	NewBlockMessage
)

func (p *peer) sendNewestBlock() {
	m := &Message{}

	chain := blockchain.BlockChain()
	block := blockchain.FindBlock(chain.NewestHash)
	blockAsJSON, err := utils.EncodeAsJSON(block)
	utils.HandleErr(err)
	m.MessageKind = SendNewestBlockMessage
	m.Payload = blockAsJSON

	mBytes := utils.ToBytes(m)
	p.inbox <- mBytes
}

func (p *peer) requestAllBlocks() {
	m := &Message{}

	m.MessageKind = RequestAllBlocksMessage
	mBytes := utils.ToBytes(m)

	p.inbox <- mBytes
}

func (p *peer) sendAllBlocks() {
	m := &Message{}
	blocks := blockchain.Blocks(blockchain.BlockChain())

	m.MessageKind = SendAllBlocksMessage
	blocksAsJSON, err := utils.EncodeAsJSON(blocks)
	utils.HandleErr(err)
	m.Payload = blocksAsJSON
	mBytes := utils.ToBytes(m)
	p.inbox <- mBytes
}

func (p *peer) NewBlock(newBlock *blockchain.Block) {
	m := &Message{}

	blockAsJSON, err := utils.EncodeAsJSON(newBlock)
	utils.HandleErr(err)

	m.MessageKind = NewBlockMessage
	m.Payload = blockAsJSON

	mBytes := utils.ToBytes(m)
	p.inbox <- mBytes
}

func (p *peer) NewTx(newTx *blockchain.Tx) {
	m := &Message{}

	newTxAsBytes, err := utils.EncodeAsJSON(newTx)
	utils.HandleErr(err)

	m.MessageKind = NewTransactionMessage
	m.Payload = newTxAsBytes
	mBytes := utils.ToBytes(m)

	p.inbox <- mBytes
}

func BroadcastMessage(kind MessageKind, payload []byte, p *peer) {
	switch kind {
	case SendNewestBlockMessage:
		block := &blockchain.Block{}
		err := utils.DecodeAsJSON(block, payload)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		myChain := blockchain.BlockChain()
		if block.Height > myChain.Height {
			fmt.Printf("I want to get all blocks of %s\n", p.Key)
			p.requestAllBlocks()
			break
		} else if block.Height < myChain.Height {
			fmt.Printf("Send my all blocks to %s\n", p.Key)
			p.sendNewestBlock()
			break
		} else {
			fmt.Printf("we are same blockchain 🤘\n")
			break
		}
	case RequestAllBlocksMessage:
		fmt.Printf("I've sending all blocks to %s\n", p.Key)
		p.sendAllBlocks()
	case SendAllBlocksMessage:
		fmt.Printf("I Received all blocks from %s\n\n", p.Key)
		blocks := []*blockchain.Block{}
		err := utils.DecodeAsJSON(&blocks, payload)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		blockchain.HandleSendAllBlocksMessage(blocks)
	case NewBlockMessage:
		block := &blockchain.Block{}
		err := utils.DecodeAsJSON(block, payload)
		utils.HandleErr(err)

		blockchain.HandleNewBlockMessage(block)
		fmt.Printf("I Received new block from %s\n", p.Key)
	case NewTransactionMessage:
		tx := &blockchain.Tx{}
		err := utils.DecodeAsJSON(tx, payload)
		utils.HandleErr(err)

		blockchain.HandleNewTxMessage(tx)
		fmt.Printf("I Received new Transactions from %s\n", p.Key)
	default:
		break
	}
}
