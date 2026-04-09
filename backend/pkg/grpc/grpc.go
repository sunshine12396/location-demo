package grpc

import "google.golang.org/grpc"

func (i *impl) Inject() grpc.ServiceRegistrar {
	return i.server
}
