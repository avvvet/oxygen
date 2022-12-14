package blockchain

import (
	"bytes"
	"encoding/json"
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

func InitChainLedger(path string) (*Chain, error) {
	ledger, err := kv.NewLedger(path)
	if err != nil {
		logger.Sugar().Warn("unable to initialize chain ledger.")
		return nil, err
	}
	return &Chain{ledger, &Block{BlockHeight: -1}}, err
}

/* checks if any block exists */
func (c *Chain) IsGenesisTx() bool {
	iter := c.Ledger.Db.NewIterator(nil, nil)

	return !iter.Last() /*true if record existis */
}

// func InitChain(txOutput *TxOutput) (*Chain, error) {
// 	ledger, err := kv.NewLedger("./ledger/store")
// 	if err != nil {
// 		logger.Sugar().Fatal("unable to initialize ledger.")
// 	}

// 	iter := ledger.Db.NewIterator(nil, nil)
// 	if !iter.Last() {
// 		tx := NatureTx(txOutput, "genesis nature token")
// 		block, err := Genesis(tx)
// 		if err != nil {
// 			logger.Sugar().Fatal("unable to create genesis block.")
// 		}
// 		b, _ := json.Marshal(block)

// 		err = ledger.Upsert([]byte(strconv.Itoa(block.BlockHeight)), b)
// 		if err != nil {
// 			logger.Sugar().Fatal("unable to store data")
// 		}
// 		iter.Release()
// 		return &Chain{ledger, block}, err
// 	}

// 	lastblock := &Block{}
// 	err = json.Unmarshal(iter.Value(), lastblock)
// 	if err != nil {
// 		logger.Sugar().Fatal("unable to get block from store")
// 	}
// 	iter.Release()
// 	return &Chain{ledger, lastblock}, err
// }

func (c *Chain) ChainBlock(data string, txs []*Transaction) (*Block, error) {
	if c.LastBlock.BlockHeight != -1 {
		lastBlock := c.LastBlock
		newblock, err := CreateBlock(data, txs, lastBlock.Hash)
		newblock.BlockHeight = lastBlock.BlockHeight + 1
		if err != nil {
			return nil, err
		}

		b, _ := json.Marshal(newblock)
		err = c.Ledger.Upsert([]byte(strconv.Itoa(newblock.BlockHeight)), b)
		if err != nil {
			return nil, err
		}
		c.LastBlock = newblock

		return newblock, nil
	}

	newblock, err := Genesis(txs)
	if err != nil {
		logger.Sugar().Warn("unable to create genesis block.")
		return nil, err
	}
	newblock.BlockHeight = c.LastBlock.BlockHeight + 1

	b, _ := json.Marshal(newblock)
	err = c.Ledger.Upsert([]byte(strconv.Itoa(newblock.BlockHeight)), b)
	if err != nil {
		return nil, err
	}
	c.LastBlock = newblock

	return newblock, nil
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
				if txOutput.IsTokenOwner(address) {
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

func (c *Chain) isSpent(address string, txid []byte, txOutputIndex int) bool {
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
			for _, txInput := range tx.Inputs {

				if bytes.Equal(txid, txInput.ID) && txInput.OutputIndex == txOutputIndex {
					return true
				}

			}
		}
		i++
	}
	return false
}
