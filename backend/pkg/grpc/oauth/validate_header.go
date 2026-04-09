package oauth

import (
	"strings"

	"github.com/example/location-demo/pkg/stderr"
)

func validateHeader(authHeader string) error {
	if strings.TrimSpace(authHeader) == "" {
		return stderr.NewUnauthorizedError("authorization header is empty")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return stderr.NewUnauthorizedError("authorization header is invalid")
	}
	return nil
}
