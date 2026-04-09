package rsa

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
)

func (i *impl) Verify(plainText, signature string) error {
	// Base64 Decode signature
	signatureBytes, err := base64.RawURLEncoding.DecodeString(signature)
	if err != nil {
		return err
	}

	// Verify signature
	hashed := sha256.Sum256([]byte(plainText))
	if err := rsa.VerifyPKCS1v15(i.publicKey, crypto.SHA256, hashed[:], signatureBytes); err != nil {
		return err
	}

	return nil
}
