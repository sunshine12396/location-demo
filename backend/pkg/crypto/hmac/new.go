package hmac

import (
	"encoding/base64"
	"errors"
)

// impl is the implementation of the Key interface.
// It holds the HMAC key used for signing and verification.
type impl struct {
	key keyBytes // The raw key bytes used for HMAC operations
}

// New creates a new HMAC key from a base64-encoded string.
// The key must be 32 bytes long when decoded.
// Returns an error if the string is not valid base64 or the key length is incorrect.
func New(base64Str string) (Key, error) {
	key, err := base64.RawURLEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}
	if len(key) != 32 {
		return nil, errors.New("key must be 32 bytes long")
	}
	return &impl{
		key: key,
	}, nil
}
