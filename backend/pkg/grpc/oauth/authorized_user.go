package oauth

import (
	"context"

	"github.com/gin-gonic/gin"
)

type IdentityUser struct {
	UUID        string
	ClientID    string
	Roles       []string
	Permissions []string
	Metadata    map[string]interface{}
}

func identityUser(c *gin.Context) IdentityUser {
	if user, ok := c.Get(KeyIdentityUser); ok {
		return user.(IdentityUser)
	}
	return IdentityUser{}
}

func attachIdentity(c *gin.Context, claims Claims) {
	c.Set(KeyIdentityUser, IdentityUser{
		UUID:        claims.GetSub(),
		ClientID:    claims.GetClientID(),
		Roles:       claims.GetRoles(),
		Permissions: claims.GetPermissions(),
		Metadata:    claims.GetMetadata(),
	})
}

func gRPCIdentityUser(ctx context.Context) IdentityUser {
	identity, ok := ctx.Value(KeyIdentityUser).(IdentityUser)
	if ok {
		return identity
	}
	return IdentityUser{}
}
