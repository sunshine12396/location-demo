package oauth

type AuthorizationServerOptions func(*authorizationImpl)

func AuthorizedClientsOption(authorizedClientMap map[string]string) AuthorizationServerOptions {
	return func(c *authorizationImpl) {
		c.authorizedClientMap = authorizedClientMap
	}
}
