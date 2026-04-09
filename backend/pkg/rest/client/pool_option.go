package client

import (
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/proxy"
)

type PoolOption func(c *http.Client, t *http.Transport)

func WithProxy(proxyURL string) PoolOption {
	return func(c *http.Client, t *http.Transport) {
		if proxyURL == "" {
			t.Proxy = nil
			return
		}
		parsedURL, err := url.Parse(proxyURL)
		if err != nil {
			log.Fatalf("failed to parse proxy URL. Err: %v", err)
			return
		}

		t.Proxy = http.ProxyURL(parsedURL)
	}
}

func WithSOCK5(address string) PoolOption {
	return func(c *http.Client, t *http.Transport) {
		dialer, err := proxy.SOCKS5("tcp", address, nil, proxy.Direct)
		if err != nil {
			log.Fatalf("failed to create socks5 dialer. Err: %v", err)
			return
		}

		dial, ok := dialer.(proxy.ContextDialer)
		if !ok {
			log.Fatalf("failed to cast socks5 dialer to ContextDialer")
			return
		}
		t.DialContext = dial.DialContext
	}
}
