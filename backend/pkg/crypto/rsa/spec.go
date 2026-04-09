package rsa

type Cipher interface {
	Encrypt(plainText string) (string, error)
	Decrypt(encryptedText string) (string, error)
	Sign(plainText string) (string, error)
	Verify(plainText, signature string) error
}

type Private interface {
	Decrypt(encryptedText string) (string, error)
	Sign(plainText string) (string, error)
}

type Public interface {
	Encrypt(plainText string) (string, error)
	Verify(plainText, signature string) error
}
