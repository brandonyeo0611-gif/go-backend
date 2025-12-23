package middleware

import (
	"net/http"
	"strings"
	"context"

	"github.com/CVWO/sample-go-app/internal/auth"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"Authorization header required"}`, 
		http.StatusUnauthorized)
			return
		}

		tokenStr := strings.Replace(authHeader, "Bearer ", "", 1) 
		// to get the data

		// to validate the token
		claims, err := auth.ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, `{"error":"Invalid token"}`, http.StatusUnauthorized)
		return

		}

		// pass user to handler via context
		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		ctx = context.WithValue(ctx, usernameKey, claims.Username)

		next.ServeHTTP(w,r.WithContext(ctx))
	})
}