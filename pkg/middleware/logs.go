package middleware

import (
	"fmt"
	"net/http"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("before")  // обработка до вызова следующего хэндлера
		next.ServeHTTP(w, req) // вызываем следующий хэндлер в цепочке
		fmt.Println("after")   // обработка после вызова следующего хэндлера
	})
}
