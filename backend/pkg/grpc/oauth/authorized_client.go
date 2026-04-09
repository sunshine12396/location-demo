package oauth

import (
	"context"

	"github.com/gin-gonic/gin"
)

type AuthorizedClient struct {
	ClientID     string
	ClientSecret string
}

func identityClient(c *gin.Context) AuthorizedClient {
	if client, ok := c.Get(KeyIdentityClient); ok {
		return client.(AuthorizedClient)
	}
	return AuthorizedClient{}
}

func attachClient(c *gin.Context, authorizedClient AuthorizedClient) {
	c.Set(KeyIdentityClient, authorizedClient)
}

func gRPCIdentityClient(ctx context.Context) AuthorizedClient {
	identity, ok := ctx.Value(KeyIdentityClient).(AuthorizedClient)
	if ok {
		return identity
	}
	return AuthorizedClient{}
}
