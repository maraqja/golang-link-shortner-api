package middleware

import (
	"context"
	"link-shortner-api/configs"
	"link-shortner-api/pkg/jwt"
	"net/http"
	"strings"
)

type key string

const (
	ContextEmailKey key = "ContextEmailKey"
)

func IsAuthed(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "No Authorization Header", http.StatusUnauthorized)
			return
		}
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization Header Format", http.StatusUnauthorized)
			return
		}

		token := headerParts[1]

		_, data := jwt.NewJWT(config.Auth.Secret).Parse(token)

		// каждый HTTP запрос содержит контекст, в котором можно хранить данные req.Context()
		ctx := context.WithValue(req.Context(), ContextEmailKey, data.Email)
		// после создания контекста со значением создаем новый запрос с этим контекстом
		req_with_enhanced_context := req.WithContext(ctx) // создаем request из существующего с новым контекстом

		next.ServeHTTP(w, req_with_enhanced_context)
	})
}
