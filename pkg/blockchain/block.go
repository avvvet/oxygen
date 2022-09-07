package blockchain

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Block struct {
	Hash       [32]byte
	Data       []byte
	PrevHash   [32]byte
	Nonce      int
	Difficulty int
}

func CreateBlock(data string, prevHash [32]byte) (*Block, error) {
	block := &Block{[32]byte{}, []byte(data), prevHash, 0, Difficulty}
	pow := NewPow(block)
	nonce, hash := pow.SignBlock()
	block.Hash = hash
	block.Nonce = nonce

	if block.IsBlockValid(nonce) {
		return block, nil
	}
	return nil, errors.New("new block could not be created ")
}

func (b *Block) IsBlockValid(nonce int) bool {
	zeros := strings.Repeat("0", Difficulty)
	concat := bytes.Join([][]byte{b.Data, b.PrevHash[:], []byte(strconv.Itoa(nonce))}, []byte{})
	hashByte := sha256.Sum256(concat)
	hashString := fmt.Sprintf("%x", hashByte)

	return hashString[:Difficulty] == zeros
}

func Genesis() (*Block, error) {
	block, err := CreateBlock("genesis adam's block ", [32]byte{})
	return block, err
}
