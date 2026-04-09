package oauth

import (
	"context"
	"log"
	"strings"

	"github.com/example/location-demo/pkg/crypto/rsa"
	"github.com/example/location-demo/pkg/stderr"
	"github.com/example/location-demo/pkg/stdlog"
	"github.com/example/location-demo/pkg/stdres"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ResourceServer interface {
	BearerTokenMiddleware() func(*gin.Context)
	GRPCBearerTokenMiddleware(string) func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error)
	GinIdentityUser(c *gin.Context) IdentityUser
	GinIdentityClient(c *gin.Context) AuthorizedClient
	GRPCIdentityUser(ctx context.Context) IdentityUser
	GRPCIdentityClient(ctx context.Context) AuthorizedClient
	GRPCBearerTokenStreamMiddleware(string) grpc.StreamServerInterceptor
	OptionalBearerTokenMiddleware() func(*gin.Context)
}

type resourceImpl struct {
	logger       stdlog.Logger
	publicCipher rsa.Public
}

func NewResource(logger stdlog.Logger, publicKeyPath string) ResourceServer {
	publicKey, err := rsa.NewPublic(publicKeyPath)
	if err != nil {
		log.Fatalf("failed to initialize cipher. Err: %v", err)
	}
	return &resourceImpl{
		logger:       logger,
		publicCipher: publicKey,
	}
}

type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

func (r *resourceImpl) GRPCBearerTokenStreamMiddleware(prefixPath string) grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		// Check path
		if !strings.Contains(info.FullMethod, prefixPath) {
			return handler(srv, ss)
		}

		ctx := ss.Context()

		// Get metadata object
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return status.Error(codes.Unauthenticated, "metadata is not provided")
		}

		// Extract value from authorization header
		authHeaderValue := md["authorization"]
		if len(authHeaderValue) == 0 {
			return status.Error(codes.Unauthenticated, "authorization header is empty")
		}

		tokenValue := authHeaderValue[0]
		if !strings.HasPrefix(tokenValue, "Bearer ") {
			return status.Error(codes.Unauthenticated, "authorization header is invalid")
		}

		tokenStr := strings.TrimPrefix(tokenValue, "Bearer ")
		claims, err := verifyToken(r.publicCipher, tokenStr)
		if err != nil {
			return stderr.NewGRPCError(err)
		}

		// Attach identity to context
		newCtx := context.WithValue(ctx, KeyIdentityUser, IdentityUser{
			UUID:     claims.GetSub(),
			ClientID: claims.GetClientID(),
		})

		wrapped := &wrappedServerStream{
			ServerStream: ss,
			ctx:          newCtx,
		}
		return handler(srv, wrapped)
	}
}

func (r *resourceImpl) BearerTokenMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if err := validateHeader(authHeader); err != nil {
			stdres.UnauthorizeError(c, err)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := verifyToken(r.publicCipher, tokenStr)
		if err != nil {
			stdres.UnauthorizeError(c, err)
			return
		}

		// Attach identity to context
		attachIdentity(c, claims)

		c.Next()
	}
}

func (r *resourceImpl) OptionalBearerTokenMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		// Ignore Authorization
		if strings.TrimSpace(authHeader) == "" {
			c.Next()
			return
		}

		if err := validateHeader(authHeader); err != nil {
			stdres.UnauthorizeError(c, err)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := verifyToken(r.publicCipher, tokenStr)
		if err != nil {
			stdres.UnauthorizeError(c, err)
			return
		}

		// Attach identity to context
		attachIdentity(c, claims)

		c.Next()
	}
}

func (r *resourceImpl) GRPCBearerTokenMiddleware(prefixPath string) func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		// Check path
		if !strings.Contains(info.FullMethod, prefixPath) {
			return handler(ctx, req)
		}

		// Get metadata object
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
		}

		// Extract value from authorization header
		authHeaderValue := md["authorization"]
		if len(authHeaderValue) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization header is empty")
		}
		tokenValue := authHeaderValue[0]
		if !strings.HasPrefix(tokenValue, "Bearer ") {
			return nil, status.Error(codes.Unauthenticated, "authorization header is invalid")
		}
		tokenStr := strings.TrimPrefix(tokenValue, "Bearer ")
		claims, err := verifyToken(r.publicCipher, tokenStr)
		if err != nil {
			return nil, stderr.NewGRPCError(err)
		}

		// Attach identity to context
		ctx = context.WithValue(ctx, KeyIdentityUser, IdentityUser{
			UUID:     claims.GetSub(),
			ClientID: claims.GetClientID(),
		})

		return handler(ctx, req)
	}
}

func (r *resourceImpl) GinIdentityUser(c *gin.Context) IdentityUser {
	return identityUser(c)
}

func (r *resourceImpl) GinIdentityClient(c *gin.Context) AuthorizedClient {
	return identityClient(c)
}

func (r *resourceImpl) GRPCIdentityUser(ctx context.Context) IdentityUser {
	return gRPCIdentityUser(ctx)
}

func (r *resourceImpl) GRPCIdentityClient(ctx context.Context) AuthorizedClient {
	return gRPCIdentityClient(ctx)
}
