package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/avvvet/oxygen/pkg/blockchain"
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewDevelopment()
)

func main() {
	chain, err := blockchain.InitChain("yellow")
	if err != nil {
		fmt.Print(err)
	}
	defer chain.Ledger.Db.Close()

	tx := chain.NewTransaction("yellow", "GREEN", 400)
	chain.ChainBlock(`data {} `+strconv.Itoa(1), []*blockchain.Transaction{tx})

	tx2 := chain.NewTransaction("yellow", "GREEN", 200)
	chain.ChainBlock(`data {} `+strconv.Itoa(2), []*blockchain.Transaction{tx2})

	tx3 := chain.NewTransaction("GREEN", "BLUE", 550)
	chain.ChainBlock(`data {} `+strconv.Itoa(2), []*blockchain.Transaction{tx3})
	// for i := 1; i < 10; i++ {
	// 	chain.ChainBlock(`data {} `+strconv.Itoa(i), []*blockchain.Transaction{tx})
	// }

	var i = 0
	for {
		data, err := chain.Ledger.Get([]byte(strconv.Itoa(i)))
		if err != nil {
			break
		} else {
			block := &blockchain.Block{}
			err = json.Unmarshal(data, block)
			if err != nil {
				logger.Sugar().Fatal("unable to get block from store")
			}

			fmt.Printf("############## BlockHeight %v ############# \n", i)
			fmt.Printf("Timestamp : %s \n", time.Unix(block.Timestamp, 0).Format(time.RFC3339))
			fmt.Printf("Block Hash: %x\n", block.Hash)
			fmt.Printf("Data: %s\n", block.Data)
			fmt.Printf("Transaction ID: %x Inputs %+v  Outputs %+v\n", block.Transaction[0].ID, block.Transaction[0].Inputs, block.Transaction[0].Outputs)
			fmt.Printf("Merkle root: %x\n", block.MerkleRoot)
			fmt.Printf("Previous Hash: %x\n", block.PrevHash)
			fmt.Printf("Difficulty: %v\n", block.Difficulty)
			fmt.Printf("Nonce: %v\n", block.Nonce)
			fmt.Printf("BlockHeight: %v\n", block.BlockHeight)
		}
		i++
	}

	//iter := chain.Ledger.Db.NewIterator(nil, nil)

	// for ok := iter.First(); ok; ok = iter.Next() {
	// 	// Remember that the contents of the returned slice should not be modified, and
	// 	// only valid until the next call to Next.
	// 	key := iter.Key()

	// 	block := &blockchain.Block{}
	// 	err = json.Unmarshal(iter.Value(), block)
	// 	if err != nil {
	// 		logger.Sugar().Fatal("unable to get block from store")
	// 	}

	// 	fmt.Printf("############## %s ############# \n", key)
	// 	fmt.Printf("Block Hash: %x\n", block.Hash)
	// 	fmt.Printf("Data: %s\n", block.Data)
	// 	fmt.Printf("Previous Hash: %x\n", block.PrevHash)
	// 	fmt.Printf("Difficulty: %v\n", block.Difficulty)
	// 	fmt.Printf("Nonce: %v\n", block.Nonce)
	// 	fmt.Printf("BlockHeight: %v\n", block.BlockHeight)

	// }
	// iter.Release()
	// err = iter.Error()

}
