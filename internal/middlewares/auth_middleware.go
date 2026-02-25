package middlewares

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/hassamk122/restapi_golang/internal/auth"
	"github.com/hassamk122/restapi_golang/internal/utils"
)

// If used a plain string key, there is risk of key collison.
// other middlewares or lib might also store a value with key claims
// in the same request context
type contextKey string

const UserClaimsKey contextKey = "claims"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondWithError(res, http.StatusUnauthorized, "No token provided")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &auth.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				utils.RespondWithError(res, http.StatusBadRequest, "Invalid token signature")
			} else {
				utils.RespondWithError(res, http.StatusBadRequest, "Invalid token")
			}
			return
		}

		if token.Valid {
			ctx := context.WithValue(req.Context(), UserClaimsKey, claims)
			reqWithCtx := req.WithContext(ctx)
			next.ServeHTTP(res, reqWithCtx)
		} else {
			utils.RespondWithError(res, http.StatusUnauthorized, "Invalid token")
		}

	})
}
