package client

import (
	"bytes"
	"io"
	"net/http"
)

func (h *HttpClient) createRequest(url string, payload Payload) (*http.Request, error) {
	// Construct url
	endpointURL := payload.constructURL(url)
	var reader io.Reader
	if len(payload.Body) > 0 {
		reader = bytes.NewBuffer(payload.Body)
	}

	// Create request
	request, err := http.NewRequest(h.method, endpointURL, reader)
	if err != nil {
		return nil, err
	}

	// Set request headers
	if h.userAgent != "" {
		request.Header.Set("User-Agent", h.userAgent)
	}
	if h.contentType != "" {
		request.Header.Set("Content-Type", h.contentType)
	}
	for k, v := range h.headers { // Resource client default headers
		request.Header.Set(k, v)
	}
	for k, v := range payload.Header { // Payload headers
		request.Header.Set(k, v)
	}

	return request, nil
}
