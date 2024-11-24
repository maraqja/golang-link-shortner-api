package auth

import (
	"encoding/json"
	"fmt"
	"link-shortner-api/configs"
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
		res := LoginResponse{
			Token: "123456",
		}
		w.Header().Add("Content-Type", "application/json") // добавляем хедер
		w.WriteHeader(200) // устанавливаем статус-код 
		json.NewEncoder(w).Encode(res) // создаем json-encoder (преобразует структуры в json формат), который будет писать напрямую в респонс

		// 2 вариант, он хуже, тк требуется промежуточный буфер
		// json, err := json.Marshal(res)
		// w.Write(json)
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