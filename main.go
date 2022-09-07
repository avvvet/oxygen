package main

import (
	"fmt"

	"github.com/avvvet/oxygen/pkg/blockchain"
)

func main() {
	chain, err := blockchain.InitChain()
	if err != nil {
		fmt.Print(err)
	}

	chain.ChainBlock("First Block after adam block")
	chain.ChainBlock("Second Block after adam block")
	chain.ChainBlock("Third Block after adam block")

	for indx, block := range chain.Blocks {
		fmt.Printf("############## %v ############# \n", indx)
		fmt.Printf("Block Hash: %x\n", block.Hash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Difficulty: %v\n", block.Difficulty)
		fmt.Printf("Nonce: %v\n", block.Nonce)
	}
}
