package auth

import (
	"fmt"
	"link-shortner-api/configs"
	"net/http"
)


type AuthHandlerDependencies struct {
	*configs.Config // раносильно записи Config *configs.Config
}

type AuthHandler struct{
		*configs.Config // раносильно записи Config *configs.Config
}

func (handler *AuthHandler) Login() http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        fmt.Println("Login")
		fmt.Printf("JWT_SECRET = %s", handler.Config.Auth.Secret)
    }
}

func (handler *AuthHandler) Register() http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        fmt.Println("Register")
    }
}

func NewAuthHandler(router *http.ServeMux, dependencies *AuthHandlerDependencies) {
    handler := &AuthHandler{
		Config: dependencies.Config,
	}
    router.Handle("POST /auth/login", handler.Login())
	router.Handle("POST /auth/register", handler.Register())
}