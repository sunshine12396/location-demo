package client

import (
	"fmt"
	"net/url"
	"strings"
)

// Payload is the request payload struct representation
type Payload struct {
	// Request body
	Body []byte
	// QueryParams contains the request/query parameters
	QueryParams url.Values
	// PathVars contains the path variables used to replace placeholders
	// wrapped with {} in Client.URL
	PathVars map[string]string
	// Header contains custom request headers that will be added to the request
	// on http call.
	// The values on this field will override Client.Headers.Values
	Header map[string]string
}

func (p *Payload) constructURL(url string) string {
	// Replace path variables
	for k, v := range p.PathVars {
		url = strings.Replace(url, fmt.Sprintf(":%s", k), v, -1)
	}
	// Add query params
	if q := p.QueryParams.Encode(); q != "" {
		sep := "?"
		if strings.Contains(url, "?") {
			sep = "&"
		}
		url = url + sep + q
	}
	return url
}
