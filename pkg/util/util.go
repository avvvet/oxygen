package util

import (
	"crypto/ecdsa"
	"crypto/x509"
)

func encode(publicKey *ecdsa.PublicKey) []byte {
	encodedByte, _ := x509.MarshalPKIXPublicKey(publicKey)
	return encodedByte
}

func decode(encodedPub []byte) *ecdsa.PublicKey {
	genericPublicKey, _ := x509.ParsePKIXPublicKey(encodedPub)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return publicKey
}
