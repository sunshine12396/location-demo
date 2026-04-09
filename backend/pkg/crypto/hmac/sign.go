package hmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// Sign creates an HMAC signature for the given plainText using SHA-256.
// It returns the signature as a base64 URL-encoded string without padding.
// This method implements the Key.Sign interface method.
func (i *impl) Sign(plainText string) string {
	h := hmac.New(sha256.New, i.key)
	h.Write([]byte(plainText))
	signature := h.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(signature)
}
