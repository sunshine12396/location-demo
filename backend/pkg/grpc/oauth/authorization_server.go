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
)

type AuthorizationServer interface {
	BearerTokenMiddleware() func(*gin.Context)
	AuthorizedClientMiddleware() func(*gin.Context)
	NewToken(Claims) Token
	GinIdentityUser(c *gin.Context) IdentityUser
	GinIdentityClient(c *gin.Context) AuthorizedClient
	GRPCIdentityUser(ctx context.Context) IdentityUser
	GRPCIdentityClient(ctx context.Context) AuthorizedClient
}

type authorizationImpl struct {
	privateCipher       rsa.Private
	publicCipher        rsa.Public
	logger              stdlog.Logger
	authorizedClientMap map[string]string
}

func NewAuthorization(logger stdlog.Logger, privateKeyPath, publicKeyPath string, options ...AuthorizationServerOptions) AuthorizationServer {
	privateCipher, err := rsa.NewPrivate(privateKeyPath)
	if err != nil {
		log.Fatalf("failed to initialize private cipher. Err: %v", err)
	}
	publicCipher, err := rsa.NewPublic(publicKeyPath)
	if err != nil {
		log.Fatalf("failed to initialize public cipher. Err: %v", err)
	}
	i := &authorizationImpl{
		privateCipher: privateCipher,
		publicCipher:  publicCipher,
		logger:        logger,
	}

	for _, option := range options {
		option(i)
	}

	return i
}

func (a *authorizationImpl) BearerTokenMiddleware() func(*gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if err := validateHeader(authHeader); err != nil {
			stdres.UnauthorizeError(c, err)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := verifyToken(a.publicCipher, tokenStr)
		if err != nil {
			stdres.UnauthorizeError(c, err)
			return
		}

		// Attach identity to context
		attachIdentity(c, claims)

		c.Next()
	}
}

func (a *authorizationImpl) NewToken(claims Claims) Token {
	return newToken(a.privateCipher, claims)
}

func (a *authorizationImpl) GinIdentityUser(c *gin.Context) IdentityUser {
	return identityUser(c)
}

func (a *authorizationImpl) GinIdentityClient(c *gin.Context) AuthorizedClient {
	return identityClient(c)
}

func (a *authorizationImpl) GRPCIdentityUser(ctx context.Context) IdentityUser {
	return gRPCIdentityUser(ctx)
}

func (a *authorizationImpl) GRPCIdentityClient(ctx context.Context) AuthorizedClient {
	return gRPCIdentityClient(ctx)
}

func (a *authorizationImpl) AuthorizedClientMiddleware() func(*gin.Context) {
	return func(c *gin.Context) {
		clientIDHeader := c.GetHeader(HeaderXClientID)
		if strings.TrimSpace(clientIDHeader) == "" {
			c.Abort()
			stdres.BadRequest(c, stderr.NewBadRequest("required_client_id", "client id is required"))
			return
		}
		clientSecretHeader := c.GetHeader(HeaderXClientSecret)
		if strings.TrimSpace(clientSecretHeader) == "" {
			c.Abort()
			stdres.BadRequest(c, stderr.NewBadRequest("required_client_secret", "client secret is required"))
			return
		}
		decryptedClientSecret, err := a.privateCipher.Decrypt(clientSecretHeader)
		if err != nil {
			c.Abort()
			stdres.BadRequest(c, stderr.NewBadRequest("invalid_client_secret", "client secret is invalid"))
			return
		}

		if len(a.authorizedClientMap) > 0 {
			if val, ok := a.authorizedClientMap[clientIDHeader]; !ok || val != decryptedClientSecret {
				c.Abort()
				stdres.BadRequest(c, stderr.NewBadRequest("invalid_client", "client is invalid"))
			}
		}

		attachClient(c, AuthorizedClient{
			ClientID:     strings.TrimSpace(clientIDHeader),
			ClientSecret: strings.TrimSpace(clientSecretHeader),
		})

		c.Next()
	}
}
