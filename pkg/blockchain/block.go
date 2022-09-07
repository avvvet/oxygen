package blockchain

import (
	"bytes"
	"crypto/sha256"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
}

func (b *Block) HashBlock() {
	concat := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(concat)
	b.Hash = hash[:]
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash}
	block.HashBlock()
	return block
}

func Genesis() *Block {
	return CreateBlock("adam's block genesis", []byte{})
}
