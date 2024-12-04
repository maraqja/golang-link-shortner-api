package middleware

import (
	"fmt"
	"link-shortner-api/configs"
	"link-shortner-api/pkg/jwt"
	"net/http"
	"strings"
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

		isValid, data := jwt.NewJWT(config.Auth.Secret).Parse(token)
		fmt.Println(isValid)
		fmt.Println(data)

		next.ServeHTTP(w, req)
	})
}
