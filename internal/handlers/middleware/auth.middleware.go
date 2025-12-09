package middleware

import (
	"WalletX/pkg/utils"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
)

type contextKey string

const userIDCtx contextKey = "userID"

func writeJSONError(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func CheckUserAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get(authorizationHeader)
		if header == "" {
			writeJSONError(w, "empty auth header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			writeJSONError(w, "invalid auth header", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		if len(token) == 0 {
			writeJSONError(w, "token is empty", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ParseToken(token)
		if err != nil {
			writeJSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDCtx, claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
