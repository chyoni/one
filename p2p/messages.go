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
)

func (p *peer) sendNewestBlock() {
	m := &Message{}

	chain := blockchain.BlockChain()
	block := blockchain.FindBlock(chain.NewestHash)
	blockAsJSON, err := utils.EncodeAsJSON(block)
	if err != nil {
		utils.HandleErr(err)
	}
	m.MessageKind = SendNewestBlockMessage
	m.Payload = blockAsJSON

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
			fmt.Printf("I want to get all blocks of %s", p.Key)
			// RequestAllBlocksMessage
			break
		} else if block.Height < myChain.Height {
			fmt.Printf("Send my all blocks to %s", p.Key)
			p.sendNewestBlock()
			break
		} else {
			fmt.Printf("we are same blockchain 🤘")
			break
		}
	default:
		break
	}
}
