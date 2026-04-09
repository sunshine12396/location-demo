package client

type requestHeader struct {
	// Other request header values
	// Note: Can be overridden on Payload
	// Default: nil
	values map[string]string
}
