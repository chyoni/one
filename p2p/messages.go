package p2p

import (
	"fmt"
	"strconv"

	"github.com/chiwon99881/one/blockchain"
	"github.com/chiwon99881/one/utils"
)

type MessageKind int

type Message struct {
	MessageKind MessageKind `json:"messageKind"`
	Payload     []byte      `json:"payload"`
}

type newPeerPayload struct {
	Addr       string `json:"addr"`
	Port       string `json:"port"`
	RemotePort string `json:"remotePort"`
}

const (
	SendNewestBlockMessage MessageKind = iota
	RequestAllBlocksMessage
	NewTransactionMessage
	SendAllBlocksMessage
	NewBlockMessage
	NewPeerMessage
)

func (p *peer) sendNewestBlock() {
	m := &Message{}

	newestHash := blockchain.GetNewestHash()
	block := blockchain.FindBlock(newestHash)
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

func (p *peer) newPeer(newPeer *peer) {
	m := &Message{}

	newPeerPayload := &newPeerPayload{
		Addr:       newPeer.Addr,
		Port:       newPeer.Port,
		RemotePort: p.Port,
	}
	payload, err := utils.EncodeAsJSON(newPeerPayload)
	utils.HandleErr(err)

	m.MessageKind = NewPeerMessage
	m.Payload = payload
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
		result := blockchain.HandleSendNewestBlockMessage(block)
		switch result {
		case blockchain.ChangeMyBlockChain:
			fmt.Printf("I want to get all blocks of %s\n", p.Key)
			p.requestAllBlocks()
		case blockchain.SendMyBlockChain:
			fmt.Printf("Send my all blocks to %s\n", p.Key)
			p.sendNewestBlock()
		case blockchain.NothingToDoAnything:
			fmt.Printf("we are same blockchain 🤘\n")
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
	case NewPeerMessage:
		newPeerPayload := &newPeerPayload{}
		err := utils.DecodeAsJSON(newPeerPayload, payload)
		utils.HandleErr(err)

		remotePort, err := strconv.Atoi(newPeerPayload.RemotePort)
		utils.HandleErr(err)
		ConnectPeer(newPeerPayload.Addr, newPeerPayload.Port, remotePort, false)
	default:
		break
	}
}
