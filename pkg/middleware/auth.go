package middleware

import (
	"api/configs"
	"api/pkg/jwt"
	"context"
	"net/http"
	"strings"
)

type key string

const (
	ContextEmailKey key = "ContextEmailKey"
)

func writeUnathed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func Auth(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")
		if !strings.HasPrefix(bearerToken, "Bearer ") {
			writeUnathed(w)
			return
		}
		token := strings.TrimPrefix(bearerToken, "Bearer ")
		isValid, data := jwt.NewJWT(config.Auth.Secret).Parse(token)

		if !isValid {
			writeUnathed(w)
			return
		}
		ctx := context.WithValue(r.Context(), ContextEmailKey, data.Email)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
