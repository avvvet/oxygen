package blockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
)

const Difficulty int = 3

type ProofOfWork struct {
	Block *Block
}

func NewPow(block *Block) *ProofOfWork {
	return &ProofOfWork{block}
}

func (pow *ProofOfWork) SignBlock() (int, [32]byte) {
	var (
		nonce = 0
		hash  [32]byte
		flag  bool
	)

	for {
		flag, hash = pow.isSigned(nonce)
		if flag {
			break
		}
		nonce++
	}

	return nonce, hash
}

func (pow *ProofOfWork) isSigned(nonce int) (bool, [32]byte) {
	zeros := strings.Repeat("0", Difficulty)
	concat := bytes.Join([][]byte{pow.Block.Data, pow.Block.PrevHash[:], []byte(strconv.Itoa(nonce))}, []byte{})
	hashByte := sha256.Sum256(concat)
	hashString := fmt.Sprintf("%x", hashByte)

	return hashString[:Difficulty] == zeros, hashByte
}
