package client

import (
	"crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func WithOAuthClient(c OAuthClient) grpc.DialOption {
	return grpc.WithPerRPCCredentials(c)
}

func InsecureConnection() grpc.DialOption {
	return grpc.WithTransportCredentials(insecure.NewCredentials())
}

func SecureConnection() grpc.DialOption {
	return grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{}))
}
