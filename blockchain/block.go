package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"
)

type Block struct {
	Data      string
	Hash      string
	PrevHash  string
	Timestamp int
}

func (b *Block) hash() {
	var blockBuffer bytes.Buffer
	enc := gob.NewEncoder(&blockBuffer)
	enc.Encode(b)
	bytes := sha256.Sum256(blockBuffer.Bytes())
	hash := fmt.Sprintf("%x", bytes)
	b.Hash = hash
}

// CreateBlock is generate new block.
func CreateBlock(data, prevHash string) *Block {
	b := &Block{
		Data:      data,
		PrevHash:  prevHash,
		Timestamp: int(time.Now().Unix()),
	}
	b.hash()
	return b
}
