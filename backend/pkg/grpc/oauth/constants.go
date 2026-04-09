package oauth

const (
	KeyIdentityUser   = "identity_user"
	KeyIdentityClient = "identity_client"
)

const (
	defaultAuthJWTExpirationDuration        = 7
	defaultAuthJWTRefreshExpirationDuration = 30
)

const (
	AlgorithmHS256 string = "HS256"
	AlgorithmES256 string = "ES256"
)

const (
	TypeJWT       string = "JWT"
	TypeBearerJWT string = "Bearer+JWT"
)

const (
	ContentTypeJSON string = "JSON"
)

const (
	HeaderXClientID     string = "x-client-id"
	HeaderXClientSecret string = "x-client-secret"
)

const (
	Issuer string = "https://tbchain.io"
)

const (
	Audience string = "https://tbchain.io"
)
