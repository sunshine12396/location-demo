package client

// HttpClientConfig holds the base config for Client
type HttpClientConfig struct {
	URL     string // The URL we need to call
	Method  string // The HTTP Method to be used
	Headers map[string]string
}

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	TokenURL     string
}
