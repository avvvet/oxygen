package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type WalletAddress struct {
	PrivateKey    *ecdsa.PrivateKey
	PublicKey     *ecdsa.PublicKey
	WalletAddress string
}
type Signature struct {
	R *big.Int
	S *big.Int
}

type RawTx struct {
	SenderPublicKey       []byte
	SenderWalletAddress   string
	SenderRandomHash      [32]byte
	ReceiverPublicKey     []byte
	ReceiverWalletAddress string
	Token                 int
}

func NewWallet() *WalletAddress {
	// 1. Creating ECDSA private key (32 bytes) public key (64 bytes)
	w := &WalletAddress{}
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	w.PrivateKey = privateKey
	w.PublicKey = &w.PrivateKey.PublicKey
	// 2. Perform SHA-256 hashing on the public key (32 bytes).
	h2 := sha256.New()
	h2.Write(w.PublicKey.X.Bytes())
	h2.Write(w.PublicKey.Y.Bytes())
	digest2 := h2.Sum(nil)
	// 3. Perform RIPEMD-160 hashing on the result of SHA-256 (20 bytes).
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)
	// 4. Add version byte in front of RIPEMD-160 hash (0x00 for Main Network).
	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest3[:])
	// 5. Perform SHA-256 hash on the extended RIPEMD-160 result.
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)
	// 6. Perform SHA-256 hash on the result of the previous SHA-256 hash.
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)
	// 7. Take the first 4 bytes of the second SHA-256 hash for checksum.
	chsum := digest6[:4]
	// 8. Add the 4 checksum bytes from 7 at the end of extended RIPEMD-160 hash from 4 (25 bytes).
	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:])
	copy(dc8[21:], chsum[:])
	// 9. Convert the result from a byte string into base58.
	address := base58.Encode(dc8)
	w.WalletAddress = address
	return w
}

// func (w *WalletAddress) PrivateKey() *ecdsa.PrivateKey {
// 	return w.privateKey
// }

func (w *WalletAddress) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.PrivateKey.D.Bytes())
}

// func (w *WalletAddress) PublicKey() *ecdsa.PublicKey {
// 	return w.publicKey
// }

func (w *WalletAddress) PublicKeyStr() string {
	return fmt.Sprintf("%064x%064x", w.PublicKey.X.Bytes(), w.PublicKey.Y.Bytes())
}

// func (w *WalletAddress) BlockchainAddress() string {
// 	return w.blockchainAddress
// }

func (rtx *RawTx) Sign(pk *ecdsa.PrivateKey) *Signature {
	m, _ := json.Marshal(rtx)
	h := sha256.Sum256([]byte(m))
	r, s, _ := ecdsa.Sign(rand.Reader, pk, h[:])
	return &Signature{r, s}
}

type Transaction struct {
	senderPrivateKey       *ecdsa.PrivateKey
	senderPublicKey        *ecdsa.PublicKey
	senderWalletAddress    string
	recipientWalletAddress string
	value                  float32
}

func NewTransaction(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey,
	sender string, recipient string, value float32) *Transaction {
	return &Transaction{privateKey, publicKey, sender, recipient, value}
}
