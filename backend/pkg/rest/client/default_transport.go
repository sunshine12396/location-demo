package client

import "net/http"

func defaultTransport() *http.Transport {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConnsPerHost = defaultMaxIdleConnsPerHost
	return transport
}
