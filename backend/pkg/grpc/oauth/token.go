package oauth

import (
	"fmt"

	"github.com/example/location-demo/pkg/crypto/rsa"
)

type Token interface {
	BearerToken() (string, error)
	RefreshToken() (string, error)
}

type tokenImpl struct {
	privateCipher rsa.Private
	header        Header
	claims        Claims
}

func newToken(privateCipher rsa.Private, claims Claims) Token {
	return &tokenImpl{
		privateCipher: privateCipher,
		header:        newHeader(),
		claims:        claims,
	}
}

func (t *tokenImpl) BearerToken() (string, error) {
	t.claims.setBearerExp()
	return generateToken(t)
}

func (t *tokenImpl) RefreshToken() (string, error) {
	t.claims.setRefreshExp()
	return generateToken(t)
}

func generateToken(t *tokenImpl) (string, error) {
	headerBase64, err := t.header.Base64Encode()
	if err != nil {
		return "", err
	}

	claimsBase64, err := t.claims.Base64Encode()
	if err != nil {
		return "", err
	}

	messageBase64 := fmt.Sprintf("%s.%s", headerBase64, claimsBase64)
	signatureBase64, err := t.privateCipher.Sign(messageBase64)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.%s", messageBase64, signatureBase64), nil
}
