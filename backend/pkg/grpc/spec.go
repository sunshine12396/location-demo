package grpc

import (
	"net"

	"google.golang.org/grpc"
)

type Server interface {
	GracefulStop()
	Serve(listener net.Listener)
	Inject() grpc.ServiceRegistrar
	RegisterRoutes(...func(Server))
	Routes()
}
