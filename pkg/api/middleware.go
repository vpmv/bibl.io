package api

import (
	"github.com/go-fuego/fuego"
	"net/http"
)

func (api *API) bearerAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth, ok, err := api.auth.AuthenticateBearer(fuego.TokenFromHeader(r))

		if err != nil {
			api.internalServerError(w, r, err, "internal server error")
			return
		}

		if !ok {
			api.unauthorized(w, r)
			return
		}

		next.ServeHTTP(w, r.WithContext(addAuthorizationContext(r.Context(), auth)))
	})
}

func (api *API) hasPermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := getAuthorization(r.Context())

			if !auth.HasPermission(permission) {
				api.unauthorized(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
