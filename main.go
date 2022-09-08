package main

import (
	"encoding/json"
	"fmt"

	"github.com/avvvet/oxygen/pkg/blockchain"
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewDevelopment()
)

func main() {
	chain, err := blockchain.InitChain()
	if err != nil {
		fmt.Print(err)
	}
	defer chain.Ledger.Db.Close()

	chain.ChainBlock("First Block after adam block")
	chain.ChainBlock("Second Block after adam block")
	chain.ChainBlock("Third Block after adam block")

	iter := chain.Ledger.Db.NewIterator(nil, nil)

	for ok := iter.Last(); ok; ok = iter.Prev() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()

		block := &blockchain.Block{}
		err = json.Unmarshal(iter.Value(), block)
		if err != nil {
			logger.Sugar().Fatal("unable to get block from store")
		}

		fmt.Printf("############## %x ############# \n", key)
		fmt.Printf("Block Hash: %x\n", block.Hash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Difficulty: %v\n", block.Difficulty)
		fmt.Printf("Nonce: %v\n", block.Nonce)
		fmt.Printf("BlockHeight: %v\n", block.BlockHeight)

	}
	iter.Release()
	err = iter.Error()
	// for indx, block := range chain.Blocks {
	// 	fmt.Printf("############## %v ############# \n", indx)
	// 	fmt.Printf("Block Hash: %x\n", block.Hash)
	// 	fmt.Printf("Data: %s\n", block.Data)
	// 	fmt.Printf("Previous Hash: %x\n", block.PrevHash)
	// 	fmt.Printf("Difficulty: %v\n", block.Difficulty)
	// 	fmt.Printf("Nonce: %v\n", block.Nonce)
	// }
}
