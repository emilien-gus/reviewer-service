package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			role, ok := tokens[token]
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "role", role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
