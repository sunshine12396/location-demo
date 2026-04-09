package client

import (
	"log"

	"google.golang.org/grpc"
)

type grpcConn struct {
	conn *grpc.ClientConn
}

func NewConn(address string, options ...grpc.DialOption) GrpcConn {
	conn, err := grpc.NewClient(
		address,
		options...,
	)
	if err != nil {
		log.Fatalf("failed to initialize grpc connection. Error: %v", err)
	}
	return &grpcConn{
		conn: conn,
	}
}

func (g *grpcConn) Close() {
	defer g.conn.Close()
}

func (g *grpcConn) Dial() *grpc.ClientConn {
	return g.conn
}
