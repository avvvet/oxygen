package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

type TxInput struct {
	ID  []byte //ref which transaction
	Out int    // which output I am to use
	Sig string // something that allow the owner to use the value
}

type TxOutput struct {
	Token  int
	PubKey string // allows the receiver to unlock the value
}

type UTXO struct {
	ID            []byte
	TxoutputIndex int
	Token         int
}

func NatureTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Nature Token to %s", to)
	}

	txinput := TxInput{[]byte{}, -1, data}
	txoutput := TxOutput{1000, to}

	tx := Transaction{nil, []TxInput{txinput}, []TxOutput{txoutput}}
	tx.GenTxId()
	return &tx
}

func (tx *Transaction) GenTxId() {
	b, err := json.Marshal(tx)
	if err != nil {
		logger.Sugar().Fatal("could not encode transaction for hashing.")
	}
	hash := sha256.Sum256(b)
	tx.ID = hash[:]
}

func (tx *Transaction) IsNatureToken() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}

func (out *TxOutput) CanUnlock(data string) bool {
	return out.PubKey == data
}

func NewTX(from, to string, amount int, c *Chain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	utxo, total_utxo := c.GetUTXO(from)

	/* sort utxo */

	sort.SliceStable(utxo, func(i, j int) bool {
		return utxo[i].Token < utxo[j].Token
	})

	if total_utxo < amount {
		logger.Sugar().Info("Error: not enought funds")
		return nil
	}

	/*
	   starting from lowest token output
	   check if transaction can be made by single output
	   else use enough multiple unspent outputs (utxo)
	*/

	var enough_utxo int = 0
	for _, u := range utxo {

		enough_utxo += u.Token

		if enough_utxo >= amount { // token is enough
			input := TxInput{u.ID, u.TxoutputIndex, from}
			inputs = append(inputs, input)

			outputs = append(outputs, TxOutput{amount, to})
			outputs = append(outputs, TxOutput{enough_utxo - amount, from}) // change to sender

			break
		}

		input := TxInput{u.ID, u.TxoutputIndex, from}
		inputs = append(inputs, input)
	}

	tx := Transaction{nil, inputs, outputs}
	tx.GenTxId()

	return &tx
}
