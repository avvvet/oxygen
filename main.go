package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"time"

	"github.com/avvvet/oxygen/pkg/blockchain"
	"github.com/avvvet/oxygen/pkg/wallet"
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewDevelopment()
)

func main() {
	// Example: this will give us a 32 byte output
	randomString, err := GenerateRandomString(32)
	if err != nil {
		// Serve an appropriately vague error to the
		// user, but log the details internally.
		panic(err)
	}

	//temp create wallet address and sign the first genesis transaction output
	adamWallet := wallet.NewWallet()
	eveWallet := wallet.NewWallet()

	senderPK := encode(adamWallet.PublicKey)
	receiverPK := encode(eveWallet.PublicKey)

	natureRawTx := &wallet.RawTx{
		SenderPublicKey:       senderPK,
		SenderWalletAddress:   adamWallet.WalletAddress,
		SenderRandomHash:      sha256.Sum256([]byte(randomString)),
		ReceiverPublicKey:     senderPK,
		ReceiverWalletAddress: adamWallet.WalletAddress,
		Token:                 900,
	}

	txout := &blockchain.TxOutput{
		RawTx:     natureRawTx,
		Signature: natureRawTx.Sign(adamWallet.PrivateKey),
	}

	chain, err := blockchain.InitChain(txout)
	if err != nil {
		fmt.Print(err)
	}
	defer chain.Ledger.Db.Close()

	rawTx1 := &wallet.RawTx{
		SenderPublicKey:       senderPK,
		SenderWalletAddress:   adamWallet.WalletAddress,
		SenderRandomHash:      sha256.Sum256([]byte(randomString)),
		Token:                 400,
		ReceiverPublicKey:     receiverPK,
		ReceiverWalletAddress: eveWallet.WalletAddress,
	}

	txout1 := &blockchain.TxOutput{
		RawTx:     rawTx1,
		Signature: rawTx1.Sign(adamWallet.PrivateKey),
	}
	tx1 := chain.NewTransaction(txout1)
	chain.ChainBlock(`data {} `+strconv.Itoa(1), []*blockchain.Transaction{tx1})

	rawTx2 := &wallet.RawTx{
		SenderPublicKey:       senderPK,
		SenderWalletAddress:   adamWallet.WalletAddress,
		SenderRandomHash:      sha256.Sum256([]byte(randomString)),
		Token:                 490,
		ReceiverPublicKey:     receiverPK,
		ReceiverWalletAddress: eveWallet.WalletAddress,
	}

	txout2 := &blockchain.TxOutput{
		RawTx:     rawTx2,
		Signature: rawTx2.Sign(adamWallet.PrivateKey),
	}
	tx2 := chain.NewTransaction(txout2)
	chain.ChainBlock(`data {} `+strconv.Itoa(1), []*blockchain.Transaction{tx2})

	// tx3 := chain.NewTransaction("GREEN", "BLUE", 550)
	// chain.ChainBlock(`data {} `+strconv.Itoa(2), []*blockchain.Transaction{tx3})

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

func StreamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

func StreamToString(stream io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.String()
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

func encode(publicKey *ecdsa.PublicKey) []byte {
	encodedByte, _ := x509.MarshalPKIXPublicKey(publicKey)
	return encodedByte
}

func decode(encodedPub []byte) *ecdsa.PublicKey {
	genericPublicKey, _ := x509.ParsePKIXPublicKey(encodedPub)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return publicKey
}
