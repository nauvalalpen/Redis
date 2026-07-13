package middleware

import (
	"context"
	"net/http"
	"strings"

	"minikatalog/internal/usecase"
)

// contextKey adalah type khusus untuk menghindari collision di context
type contextKey string

const ContextUserID contextKey = "userID"

// AuthMiddleware memvalidasi token dari Header Authorization
// Format: "Authorization: Bearer <token>"
func AuthMiddleware(authUC *usecase.AuthUsecase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"Authorization header diperlukan"}`, http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, `{"error":"Format Authorization tidak valid. Gunakan: Bearer <token>"}`, http.StatusUnauthorized)
				return
			}

			token := parts[1]
			userID, err := authUC.ValidateSession(token)
			if err != nil {
				http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusUnauthorized)
				return
			}

			// Simpan userID ke context untuk digunakan handler
			ctx := context.WithValue(r.Context(), ContextUserID, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext mengambil userID dari context request
func GetUserIDFromContext(r *http.Request) string {
	val := r.Context().Value(ContextUserID)
	if val == nil {
		return ""
	}
	return val.(string)
}
