package cryptography

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/curve25519"
)

type KeyPair struct {
	PrivateKey string
	PublicKey  string
}

// https://github.com/XTLS/Xray-core/blob/main/main/commands/all/curve25519.go
func GenerateCurve25519Keys() (*KeyPair, error) {
	var err error
	var privateKey, publicKey []byte
	encoding := base64.RawURLEncoding

	privateKey = make([]byte, curve25519.ScalarSize)
	if _, err = rand.Read(privateKey); err != nil {
		return nil, err
	}

	// Modify random bytes using algorithm described at:
	// https://cr.yp.to/ecdh.html.
	privateKey[0] &= 248
	privateKey[31] &= 127 | 64

	if publicKey, err = curve25519.X25519(privateKey, curve25519.Basepoint); err != nil {
		return nil, err
	}

	return &KeyPair{
		PrivateKey: encoding.EncodeToString(privateKey),
		PublicKey:  encoding.EncodeToString(publicKey),
	}, nil
}
