package blockchain

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	Timestamp   int64
	Hash        [32]byte
	Data        []byte
	Transaction []*Transaction
	MerkleRoot  []byte
	PrevHash    [32]byte
	Nonce       int
	BlockHeight int
	Difficulty  int
}

func CreateBlock(data string, txs []*Transaction, prevHash [32]byte) (*Block, error) {
	block := &Block{
		time.Now().Unix(),
		[32]byte{},
		[]byte(data),
		txs,
		[]byte{},
		prevHash,
		0,
		0,
		Difficulty,
	}

	pow := NewPow(block)
	nonce, hash := pow.SignBlock()
	block.Hash = hash
	block.Nonce = nonce
	block.MerkleRoot = GenerateMerkleRoot(txs)

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

func Genesis(tx []*Transaction) (*Block, error) {
	block, err := CreateBlock("genesis adam's block ", tx, [32]byte{})
	return block, err
}
