package hmac

import "errors"

// Verify checks if the originalSignature is valid for the given plainText.
// It works by generating a new signature for the plainText and comparing it
// with the originalSignature. If they match, the signature is valid.
// 
// Parameters:
//   - plainText: The original text that was signed
//   - originalSignature: The signature to verify, as a base64-encoded string
//
// Returns:
//   - nil if the signature is valid
//   - an error if the signature is invalid
func (i *impl) Verify(plainText, originalSignature string) error {
	// Generate a new signature and compare with the original
	if signature := i.Sign(plainText); signature != originalSignature {
		return errors.New("verify fail")
	}

	return nil
}
