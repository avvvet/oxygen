package blockchain

import (
	"bytes"
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

func InitChain(address string) (*Chain, error) {
	ledger, err := kv.NewLedger("./ledger/store")
	if err != nil {
		logger.Sugar().Fatal("unable to initialize ledger.")
	}

	iter := ledger.Db.NewIterator(nil, nil)
	if !iter.Last() {
		tx := NatureTx(address, "genesis nature token")
		block, err := Genesis(tx)
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

func (c *Chain) ChainBlock(data string, txs []*Transaction) {
	lastBlock := c.LastBlock
	newblock, err := CreateBlock(data, txs, lastBlock.Hash)
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

func (c *Chain) GetUTXO(address string) ([]UTXO, int) {
	var utxo []UTXO
	var total_utxo int = 0
	var i = 0
	for {
		rawBlock, err := c.Ledger.Get([]byte(strconv.Itoa(i)))
		if err != nil {
			break
		}

		block := &Block{}
		err = json.Unmarshal(rawBlock, block)
		if err != nil {
			logger.Sugar().Fatal("unable to get block from store")
		}

		for _, tx := range block.Transaction {
			for txOutputIndex, txOutput := range tx.Outputs {
				if txOutput.CanUnlock(address) {
					/*
					  txoutput found , check if it is not spent
					  send this tx id , and check if it is not used/ref in any input transactions
					  also if it is ot spent accumulate token
					*/
					if !c.isSpent(address, tx.ID, txOutputIndex) {
						utxo = append(utxo, UTXO{tx.ID, txOutputIndex, txOutput.Token})
						total_utxo = total_utxo + txOutput.Token
					}
				}
			}
		}
		i++
	}

	return utxo, total_utxo
}

func (c *Chain) isSpent(address string, txid []byte, index int) bool {
	var i = 0
	for {
		rawBlock, err := c.Ledger.Get([]byte(strconv.Itoa(i)))
		if err != nil {
			break
		}

		block := &Block{}
		err = json.Unmarshal(rawBlock, block)
		if err != nil {
			logger.Sugar().Fatal("unable to get block from store")
		}

		for txOutputIndex, tx := range block.Transaction {
			for _, txInput := range tx.Inputs {
				if txInput.CanUnlock(address) {
					if bytes.Equal(txid, txInput.ID) && txOutputIndex == index {
						return true
					}
				}
			}
		}
		i++
	}
	return false
}
