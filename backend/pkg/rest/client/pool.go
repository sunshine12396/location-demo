package client

import "net/http"

// SharedPool is a custom wrapper around http.Client that sets the client timeout to zero so that it can be
// controlled by Client instead.
type SharedPool struct {
	*http.Client
}

// NewSharedPool returns a new http.Client instance with customizable options, ensures
// the efficient reuse of resources by pooling HTTP connections, thereby reducing
// overhead and improving performance across multiple clients
//
// Example:
//
//	func main() {
//		sharedPool := NewSharedPool()
//
//		thirdPartyServiceClient1 := thirdPartyService.NewClient(sharedPool)
//		thirdPartyServiceClient2 := thirdPartyService.NewClient(sharedPool)
//	}
//
//	// Inside third party service
//	func (srv thirdPartyService) Send(ctx context.Context) error {
//		srv.sharedPool.Do(ctx,... ) // Implement your logic
//	}
//
// Refer https://www.loginradius.com/blog/engineering/tune-the-go-http-client-for-high-performance/
// NewSharedPool returns a new custom http.Client instance with custom retry and timeout options based on the arguments
func NewSharedPool(opts ...PoolOption) *SharedPool {
	return &SharedPool{defaultClient(defaultTransport(), opts...)}
}
