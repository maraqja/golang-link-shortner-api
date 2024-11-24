package auth

import (
	"fmt"
	"net/http"
)


type AuthHandler struct{}

func (h *AuthHandler) Login() http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        fmt.Println("Login")
    }
}

func (h *AuthHandler) Register() http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        fmt.Println("Register")
    }
}

func NewAuthHandler(router *http.ServeMux) {
    handler := &AuthHandler{}
    router.Handle("POST /auth/login", handler.Login())
	router.Handle("POST /auth/register", handler.Register())
}