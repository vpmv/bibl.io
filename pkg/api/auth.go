package api

import (
	"context"
	"github.com/vpmv/bibl.io/pkg/dto"
)

const AuthContextKey = `Authorization`

type Authenticator interface {
	AuthenticateBearer(bearer string) (*dto.Authorization, bool, error)
}

// getAuthorization retrieves authorization DTO from context
func getAuthorization(ctx context.Context) *dto.Authorization {
	return ctx.Value(AuthContextKey).(*dto.Authorization)
}

// addAuthorizationContext sets authorization DTO
func addAuthorizationContext(ctx context.Context, auth *dto.Authorization) context.Context {
	return context.WithValue(ctx, AuthContextKey, auth)
}
