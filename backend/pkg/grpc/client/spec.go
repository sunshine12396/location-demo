package client

import (
	"context"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
)

type GrpcConn interface {
	Dial() *grpc.ClientConn
	Close()
}

type OAuthClient interface {
	Token() (*oauth2.Token, error)
	GetRequestMetadata(context.Context, ...string) (map[string]string, error)
	RequireTransportSecurity() bool
}
