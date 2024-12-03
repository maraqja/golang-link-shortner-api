package auth

import (
	"fmt"
	"link-shortner-api/configs"
	"link-shortner-api/pkg/request"
	"link-shortner-api/pkg/response"
	"net/http"
)

type AuthHandlerDependencies struct { // структура для передачи зависимостей в конструктор
	*configs.Config // раносильно записи Config *configs.Config
	*AuthService
}

type AuthHandler struct { // структура для хранения зависимостей
	*configs.Config // раносильно записи Config *configs.Config
	*AuthService
}

func NewAuthHandler(router *http.ServeMux, dependencies *AuthHandlerDependencies) {
	handler := &AuthHandler{
		Config:      dependencies.Config,
		AuthService: dependencies.AuthService,
	}
	router.Handle("POST /auth/login", handler.Login())
	router.Handle("POST /auth/register", handler.Register())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		payload, err := request.HandleBody[LoginRequest](&w, req) // в квадратные скобки передаем дженерик

		if err != nil {
			return // HandleBody сама выдаст соответствующий респонс при ошибке
		}
		fmt.Println(payload)
		data := LoginResponse{
			Token: "123456",
		}

		response.ReturnJSON(w, http.StatusOK, data)

	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		payload, err := request.HandleBody[RegisterRequest](&w, req) // в квадратные скобки передаем дженерик

		if err != nil {
			return
		}

		handler.AuthService.Register(payload.Email, payload.Password, payload.Name)
		// response.ReturnJSON(w, http.StatusOK, data)
	}
}
