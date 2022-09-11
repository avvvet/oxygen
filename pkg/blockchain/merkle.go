package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"log"

	"github.com/avvvet/merkletree"
)

type LeafContent struct {
	content string
}

//CalculateHash hashes the values of a LeafContent
func (t LeafContent) GenerateHash() ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(t.content)); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

//Equals tests for equality of two Contents
func (t LeafContent) EqualsToHash(other merkletree.Content) (bool, error) {
	return t.content == other.(LeafContent).content, nil
}

func GenerateMerkleRoot(txs []*Transaction) []byte {
	//Build list of Content to build tree
	var list []merkletree.Content

	for _, tx := range txs {
		h := tx.HashTx()
		list = append(list, LeafContent{content: hex.EncodeToString(h[:])})
	}

	//Create a new Merkle Tree from the list of Content
	t, err := merkletree.NewTree(list)
	if err != nil {
		log.Fatal(err)
	}

	//Merkle Root of the tree
	return t.MerkleRoot()
}
