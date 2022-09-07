package main

import (
	"fmt"

	"github.com/avvvet/oxygen/pkg/blockchain"
)

func main() {
	chain := blockchain.InitChain()

	chain.ChainBlock("First Block after adam block")
	chain.ChainBlock("Second Block after adam block")
	chain.ChainBlock("Third Block after adam block")

	for indx, block := range chain.Blocks {
		fmt.Printf("############## %v ############# \n", indx)
		fmt.Printf("Block Hash: %x\n", block.Hash)
		fmt.Printf("Previous Hash: %s\n", block.Data)
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
	}

}
