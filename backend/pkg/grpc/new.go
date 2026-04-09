package grpc

import (
	"io"

	"github.com/example/location-demo/pkg/enum"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Config struct {
	Environment        string
	CertPath           string
	KeyPath            string
	writer             io.Writer
	interceptors       []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
}

type impl struct {
	server *grpc.Server
	routes []func(Server)
}

func New(cfg Config, opts ...Option) Server {
	// Options
	for _, opt := range opts {
		opt(&cfg)
	}

	// Create server
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(cfg.interceptors...))

	if cfg.Environment != enum.EnvProd.String() {
		// Register reflection service on gRPC server.
		reflection.Register(server)
	}

	return &impl{
		server: server,
	}
}
