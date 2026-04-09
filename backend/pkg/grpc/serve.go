package grpc

import (
	"log"
	"net"
)

func (i *impl) Serve(listener net.Listener) {
	log.Fatal(i.server.Serve(listener))
}
