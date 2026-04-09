package oauth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/example/location-demo/pkg/crypto/rsa"
	"github.com/example/location-demo/pkg/stderr"
)

func verifyToken(publicCipher rsa.Public, tokenStr string) (Claims, error) {
	// Reason: error message DOESN'T HAVE TO CLEAN for preventing guessing logic
	const generalErrorMessage = "token is invalid"

	// Split token into header, payload and signature
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, stderr.NewUnauthorizedError(generalErrorMessage)
	}

	// Combine header and payload to a message
	messageBase64 := fmt.Sprintf("%s.%s", parts[0], parts[1])
	if err := publicCipher.Verify(messageBase64, parts[2]); err != nil {
		return nil, stderr.NewUnauthorizedError(generalErrorMessage)
	}

	// Decode payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, stderr.NewUnauthorizedError(generalErrorMessage)
	}

	// Unmarshal payload to claims
	claims := NewEmptyClaims()
	if err := json.Unmarshal(payloadBytes, claims); err != nil {
		return nil, stderr.NewUnauthorizedError(generalErrorMessage)
	}

	// Check if a token is expired
	if claims.IsExpired() {
		return nil, stderr.NewUnauthorizedError("token is expired")
	}

	return claims, nil
}
