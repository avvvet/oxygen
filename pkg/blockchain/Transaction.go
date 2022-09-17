package blockchain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"math/big"
	"sort"

	"github.com/avvvet/oxygen/pkg/wallet"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

type TxInput struct {
	ID          []byte //ref which transaction
	OutputIndex int    // which output I am to use
}

type TxOutput struct {
	RawTx      *wallet.RawTx
	Signature  *wallet.Signature
	TokenOwner string
	Token      int
}

type UTXO struct {
	ID            []byte
	TxoutputIndex int
	Token         int
}

type Signature struct {
	R *big.Int
	S *big.Int
}

func NatureTx(txOutput *TxOutput, data string) *Transaction {
	// if data == "" {
	// 	data = fmt.Sprintf("Nature Token to %s", txOutput.RawTx.ReceiverWalletAddress)
	// }

	txinput := TxInput{[]byte{}, -1}

	txOutput.Token = txOutput.RawTx.Token
	txOutput.TokenOwner = txOutput.RawTx.ReceiverWalletAddress
	tx := Transaction{nil, []TxInput{txinput}, []TxOutput{*txOutput}}
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
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].OutputIndex == -1
}

func (in *TxInput) CanUnlock(data string) bool {
	return true
}

func (out *TxOutput) IsTokenOwner(tokenOwner string) bool {
	// 1 Verify rawTx signature by rawTx.Sender PublicKey, when true it means the rawTx is valid and signed by the sender

	// 2 check if TokenOwner is equal to rawTx receiver or sender address

	if out.VerifyRawTxSignature() && (out.TokenOwner == out.RawTx.ReceiverWalletAddress || out.TokenOwner == out.RawTx.SenderWalletAddress) && tokenOwner == out.TokenOwner {
		return true
	}

	//3 Check Signature is never used in the oxygen blockchain (this avoid duble valid spent)

	return false
}

func (c *Chain) NewTransaction(txout *TxOutput) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	/* verify rawTx signature is valid */
	if !txout.VerifyRawTxSignature() {
		return nil
	}

	utxo, total_utxo := c.GetUTXO(txout.RawTx.SenderWalletAddress)

	/* sort utxo */

	sort.SliceStable(utxo, func(i, j int) bool {
		return utxo[i].Token < utxo[j].Token
	})

	if total_utxo < txout.RawTx.Token {
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

		if enough_utxo >= txout.RawTx.Token { // token is enough
			input := TxInput{u.ID, u.TxoutputIndex}
			inputs = append(inputs, input)

			/* output for receiver */
			txout.Token = txout.RawTx.Token
			txout.TokenOwner = txout.RawTx.ReceiverWalletAddress
			outputs = append(outputs, *txout)

			/*change for sender*/
			txout.Token = enough_utxo - txout.RawTx.Token
			txout.TokenOwner = txout.RawTx.SenderWalletAddress
			outputs = append(outputs, *txout)

			break
		}

		input := TxInput{u.ID, u.TxoutputIndex}
		inputs = append(inputs, input)
	}

	tx := Transaction{nil, inputs, outputs}
	tx.GenTxId()

	return &tx
}

func (tx *Transaction) HashTx() [32]byte {
	b, err := json.Marshal(tx)
	if err != nil {
		logger.Sugar().Fatal("transaction encoding error")
	}

	return sha256.Sum256(b)
}

func (out *TxOutput) VerifyRawTxSignature() bool {
	rawtx, _ := json.Marshal(out.RawTx)
	h := sha256.Sum256([]byte(rawtx))

	SenderPK := Decode(out.RawTx.SenderPublicKey)
	flag := ecdsa.Verify(SenderPK, h[:], out.Signature.R, out.Signature.S)

	return flag
}

func Decode(encodedPub []byte) *ecdsa.PublicKey {
	genericPublicKey, _ := x509.ParsePKIXPublicKey(encodedPub)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return publicKey
}
