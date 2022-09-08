package blockchain

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/avvvet/oxygen/pkg/kv"
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewDevelopment()
)

type Chain struct {
	Ledger    *kv.Ledger
	LastBlock *Block
}

func InitChain() (*Chain, error) {
	ledger, err := kv.NewLedger()
	if err != nil {
		logger.Sugar().Fatal("unable to initialize global state db.")
	}

	iter := ledger.Db.NewIterator(nil, nil)
	if !iter.Last() {
		block, err := Genesis()
		if err != nil {
			logger.Sugar().Fatal("unable to create genesis block.")
		}
		b, _ := json.Marshal(block)

		err = ledger.Upsert([]byte(strconv.Itoa(block.BlockHeight)), b)
		if err != nil {
			logger.Sugar().Fatal("unable to store data")
		}
		iter.Release()
		return &Chain{ledger, block}, err
	}

	lastblock := &Block{}
	err = json.Unmarshal(iter.Value(), lastblock)
	if err != nil {
		logger.Sugar().Fatal("unable to get block from store")
	}
	iter.Release()
	return &Chain{ledger, lastblock}, err
}

func (c *Chain) ChainBlock(data string) {
	lastBlock := c.LastBlock
	newblock, err := CreateBlock(data, lastBlock.Hash)
	newblock.BlockHeight = lastBlock.BlockHeight + 1
	if err != nil {
		fmt.Print(err)
	} else {
		b, _ := json.Marshal(newblock)
		err = c.Ledger.Upsert([]byte(strconv.Itoa(newblock.BlockHeight)), b)
		if err != nil {
			logger.Sugar().Fatal("unable to store data")
		}
		c.LastBlock = newblock
	}
}
