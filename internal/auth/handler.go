package auth

import (
	"encoding/json"
	"fmt"
	"link-shortner-api/configs"
	"link-shortner-api/pkg/response"
	"net/http"
)

type AuthHandlerDependencies struct { // структура для передачи зависимостей в конструктор
	*configs.Config // раносильно записи Config *configs.Config
}

type AuthHandler struct { // структура для хранения зависимостей
	*configs.Config // раносильно записи Config *configs.Config
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var payload LoginRequest
		err := json.NewDecoder(req.Body).Decode(&payload) // декодируем боди в структуру
		if err != nil {
			/* будет ошибка если передать:
			1. невалидный json
			2. несоотуетствуююий тип данных
			3. несоответствующий тип данных у полей

			НО если просто не передать поля структуры, то ошибки не будет:
			В этом случае полям структуры будет присвоены пустые значия (zero value):
				для string - ""
				для int - 0
				для bool - false
				для вложенных структур - структуру с пустыми значениями
				для указателей (вложенные структуры можно определять указателями) - nil
			*/
			response.ReturnJSON(w, http.StatusBadRequest, error.Error(err))
			return
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
