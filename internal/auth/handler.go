package auth

import (
	"fmt"
	"link-shortner-api/configs"
	"link-shortner-api/pkg/response"
	"net/http"
)


type AuthHandlerDependencies struct { // структура для передачи зависимостей в конструктор
	*configs.Config // раносильно записи Config *configs.Config
}

type AuthHandler struct{ // структура для хранения зависимостей
		*configs.Config // раносильно записи Config *configs.Config
}

func (handler *AuthHandler) Login() http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        fmt.Println("Login")
		fmt.Printf("JWT_SECRET = %s", handler.Config.Auth.Secret)
		data := LoginResponse{
			Token: "123456",
		}

		response.ReturnJSON(w, http.StatusOK, data)

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