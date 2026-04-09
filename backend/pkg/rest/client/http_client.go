package client

import (
	"context"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// A HttpClient describes an HTTP endpoint's client. This client is
// mainly used to send HTTP request based on the request configuration set into
// the client and the response handling logic configured within the client.
type HttpClient struct {
	// HTTP client to be used to execute HTTP call.
	pool *SharedPool

	// URL for this Client. This value can contain path
	// variable placeholders for substitution when sending the HTTP request
	url string

	// HTTP method
	method string

	// Name of the service to call. Together with serviceName will form the label in logger fields & error code prefixes.
	serviceName string

	// User agent value (i.e. RFC7231)
	userAgent string

	// Content MIME type
	// Default: application/json
	contentType string

	// Default request header configuration
	headers map[string]string

	timeoutAndRetryOption timeoutAndRetryOption

	// Disable request body logging
	// Default: false,
	disableReqBodyLogging bool

	// Disable response body logging
	// Default: false,
	disableRespBodyLogging bool

	// Disable request/response log redaction
	// Default: false,
	disableLogRedaction bool
}

func newHttpClient(cfg HttpClientConfig, pool *SharedPool, opts ...HttpClientOption) (*HttpClient, error) {
	client := &HttpClient{
		pool: pool,
		timeoutAndRetryOption: timeoutAndRetryOption{
			maxRetries:         defaultMaxRetries,
			maxWaitPerTry:      defaultMaxWaitPerTry * time.Second,
			maxWaitInclRetries: defaultMaxWaitInclRetries * time.Second,
			retryOnTimeout:     defaultRetryOnTimeout,
			retryOnStatusCodes: make(map[int]bool),
		},
		contentType: defaultContentType,
	}

	client.url = strings.TrimSpace(cfg.URL)
	if client.url == "" {
		return nil, ErrMissingURL
	}
	client.method = strings.TrimSpace(cfg.Method)
	if client.method == "" {
		return nil, ErrMissingMethod
	}
	client.headers = cfg.Headers

	for _, opt := range opts {
		opt(client)
	}

	if err := client.timeoutAndRetryOption.validate(); err != nil {
		return nil, err
	}
	return client, nil
}

// NewHttpClient returns a new Client instance based on the arguments
func NewHttpClient(cfg HttpClientConfig, pool *SharedPool, opts ...HttpClientOption) (*HttpClient, error) {
	return newHttpClient(cfg, pool, opts...)
}

func NewOAuthHttpClient(cfg HttpClientConfig, oauthCfg OAuthConfig, pool *SharedPool, opts ...HttpClientOption) (*HttpClient, error) {
	c, err := newHttpClient(cfg, pool, opts...)
	if err != nil {
		return nil, err
	}
	conf := &clientcredentials.Config{
		ClientID:     oauthCfg.ClientID,
		ClientSecret: oauthCfg.ClientSecret,
		TokenURL:     oauthCfg.TokenURL,
	}
	c.pool.Client = conf.Client(context.WithValue(context.Background(), oauth2.HTTPClient, c.pool.Client))
	return c, nil
}
