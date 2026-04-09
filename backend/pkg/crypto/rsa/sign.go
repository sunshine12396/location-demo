package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

func (i *impl) Sign(plainText string) (string, error) {
	hashed := sha256.Sum256([]byte(plainText))
	signature, err := i.privateKey.Sign(rand.Reader, hashed[:], crypto.SHA256)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(signature), nil
}
