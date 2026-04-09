package hmac

// Key is the interface for HMAC cryptographic operations.
// It provides methods for signing data and verifying signatures.
type Key interface {
	// Sign creates an HMAC signature for the given plainText.
	// Returns the signature as a base64-encoded string.
	Sign(plainText string) string

	// Verify checks if the originalSignature is valid for the given plainText.
	// Returns nil if the signature is valid, or an error if it's invalid.
	Verify(plainText, originalSignature string) error
}
