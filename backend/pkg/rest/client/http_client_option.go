package client

// HttpClientOption alters the behavior of the HttpClient
type HttpClientOption func(*HttpClient)

// WithServiceName set the serviceName into HttpClient
func WithServiceName(serviceName string) HttpClientOption {
	return func(h *HttpClient) {
		h.serviceName = serviceName
	}
}
