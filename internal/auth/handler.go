package auth

import (
	"link-shortner-api/configs"
	"link-shortner-api/pkg/jwt"
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

		email, err := handler.AuthService.Login(payload.Email, payload.Password)
		if err != nil {
			if err.Error() == ErrWrongCredentials {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		token, err := jwt.NewJWT(handler.Config.Auth.Secret).Create(jwt.JWTData{
			Email: email,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := LoginResponse{
			Token: token,
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

		email, err := handler.AuthService.Register(payload.Email, payload.Password, payload.Name)
		if err != nil {
			if err.Error() == ErrUserExists {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := jwt.NewJWT(handler.Config.Auth.Secret).Create(jwt.JWTData{
			Email: email,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := RegisterResponse{
			Token: token,
		}

		response.ReturnJSON(w, http.StatusOK, data)
	}
}
