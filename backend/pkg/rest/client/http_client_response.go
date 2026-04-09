package client

import "net/http"

// Response is the result of the http call
type Response struct {
	Status int
	Body   []byte
	Header http.Header
}
