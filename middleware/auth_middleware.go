package middleware

import (
	//"ecommerce/utils"
	"net/http"
	"strings"
)

type TokenVerifier interface {
	VerifyToken(tokenString string) (string, error)
}

func Auth(verifier TokenVerifier, next http.Handler) http.Handler {
	// next http.Handler: next HTTP handler to call

	// Returns an http.Handler that wraps next with authentication logic
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized - Missing Token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ") // Removes "Bearer " from the header

		_, err := verifier.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r) // passes the request to next, allowing the protected route to execute
	})
}
