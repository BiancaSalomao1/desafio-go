package middleware

import (
	"context"
	"net/http"
	"strings"

	"orders-api/internal/security"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
	EmailKey  contextKey = "email"
)

func Auth(
	jwtSecret string,
	tokenStore security.TokenStore,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			header := r.Header.Get("Authorization")

			if header == "" {

				http.Error(
					w,
					"authorization header required",
					http.StatusUnauthorized,
				)

				return
			}

			if !strings.HasPrefix(header, "Bearer ") {

				http.Error(
					w,
					"invalid authorization header",
					http.StatusUnauthorized,
				)

				return
			}

			token := strings.TrimPrefix(
				header,
				"Bearer ",
			)

			claims, err := security.ValidateToken(
				token,
				jwtSecret,
			)

			if err != nil {

				http.Error(
					w,
					"invalid token",
					http.StatusUnauthorized,
				)

				return
			}

			exists, err := tokenStore.Exists(
				r.Context(),
				token,
			)

			if err != nil {

				http.Error(
					w,
					"failed to validate token",
					http.StatusInternalServerError,
				)

				return
			}

			if !exists {

				http.Error(
					w,
					"token is no longer active",
					http.StatusUnauthorized,
				)

				return
			}

			ctx := context.WithValue(
				r.Context(),
				UserIDKey,
				claims["sub"],
			)

			ctx = context.WithValue(
				ctx,
				EmailKey,
				claims["email"],
			)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		})
	}
}
