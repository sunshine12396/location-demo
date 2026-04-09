package client

import (
	"net/http"
	"time"
)

func defaultClient(transport *http.Transport, opts ...PoolOption) *http.Client {
	// Default client
	client := &http.Client{
		Timeout: defaultClientTimeout * time.Second,
	}

	// Apply options
	for _, opt := range opts {
		opt(client, transport)
	}

	// Set transport
	client.Transport = transport

	return client
}
