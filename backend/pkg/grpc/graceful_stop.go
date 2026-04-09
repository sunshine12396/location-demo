package grpc

func (i *impl) GracefulStop() {
	i.server.GracefulStop()
}
