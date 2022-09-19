package util

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"io"
	"math/big"
)

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

func EncodePrivateKey(privateKey *ecdsa.PrivateKey) []byte {
	encodedByte, _ := x509.MarshalECPrivateKey(privateKey)
	return encodedByte
}

func DecodePrivateKey(encodedPrivate []byte) *ecdsa.PrivateKey {
	privateKey, _ := x509.ParseECPrivateKey(encodedPrivate) //x509
	return privateKey
}

func EncodePublicKey(publicKey *ecdsa.PublicKey) []byte {
	encodedByte, _ := x509.MarshalPKIXPublicKey(publicKey)
	return encodedByte
}

func DecodePublicKey(encodedPub []byte) *ecdsa.PublicKey {
	genericPublicKey, _ := x509.ParsePKIXPublicKey(encodedPub)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return publicKey
}
