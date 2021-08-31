package p2p

import (
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
)

func (p *peer) sendAllBlock() {
	m := Message{}

	chain := blockchain.BlockChain()
	block := blockchain.FindBlock(chain.NewestHash)
	blockAsJSON, err := utils.EncodeAsJSON(block)
	if err != nil {
		utils.HandleErr(err)
	}
	m.MessageKind = SendNewestBlockMessage
	m.Payload = blockAsJSON

	p.inbox <- m
}
