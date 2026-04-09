package oauth

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/example/location-demo/pkg/utils/datetime"
)

type Claims interface {
	Base64Encode() (string, error)
	IsExpired() bool

	GetSub() string
	GetClientID() string
	GetRoles() []string
	GetPermissions() []string
	GetMetadata() map[string]interface{}
	GetAudience() string
	GetIssuer() string

	setBearerExp()
	setRefreshExp()
}

type claimsImpl struct {
	Iss         string                 `json:"iss"` // issuer
	Sub         string                 `json:"sub"` // user id
	Cid         string                 `json:"cid"` // client id
	Iat         int64                  `json:"iat"` // issued at
	Aud         string                 `json:"aud"` // audience
	Nbf         int64                  `json:"nbf"` // not before
	Exp         int64                  `json:"exp"` // expiration time
	Roles       []string               `json:"roles,omitempty"`
	Permissions []string               `json:"permissions,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`

	bearerExp  int64 // bearer expiration time
	refreshExp int64 // refresh expiration time
}

func NewClaims(sub, cid string, bearerExp, refreshExp int) Claims {
	now := time.Now().UTC()
	return &claimsImpl{
		Iss: Issuer,
		Sub: sub,
		Cid: cid,
		Aud: Audience,
		Iat: now.Unix(),
		Nbf: now.Unix(),

		bearerExp:  now.Add(time.Duration(bearerExp) * datetime.DurationDay).Unix(),
		refreshExp: now.Add(time.Duration(refreshExp) * datetime.DurationDay).Unix(),
	}
}

func NewRBACClaim(sub, cid string, roles, permissions []string, metadata map[string]interface{}, bearerExp, refreshExp int) Claims {
	now := time.Now().UTC()
	return &claimsImpl{
		Iss:         Issuer,
		Sub:         sub,
		Cid:         cid,
		Aud:         Audience,
		Iat:         now.Unix(),
		Nbf:         now.Unix(),
		Roles:       roles,
		Permissions: permissions,
		Metadata:    metadata,

		bearerExp:  now.Add(time.Duration(bearerExp) * datetime.DurationDay).Unix(),
		refreshExp: now.Add(time.Duration(refreshExp) * datetime.DurationDay).Unix(),
	}
}

func NewEmptyClaims() Claims {
	return &claimsImpl{}
}

func (c *claimsImpl) Base64Encode() (string, error) {
	claimsBytes, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(claimsBytes), nil
}

func (c *claimsImpl) IsExpired() bool {
	now := time.Now().UTC()
	return now.Unix() > c.Exp
}

func (c *claimsImpl) GetSub() string {
	return c.Sub
}

func (c *claimsImpl) GetClientID() string {
	return c.Cid
}

func (c *claimsImpl) setBearerExp() {
	c.Exp = c.bearerExp
}

func (c *claimsImpl) setRefreshExp() {
	c.Exp = c.refreshExp
}

func (c *claimsImpl) GetRoles() []string {
	return c.Roles
}

func (c *claimsImpl) GetPermissions() []string {
	return c.Permissions
}

func (c *claimsImpl) GetMetadata() map[string]interface{} {
	return c.Metadata
}

func (c *claimsImpl) GetAudience() string {
	return c.Aud
}

func (c *claimsImpl) GetIssuer() string {
	return c.Iss
}
