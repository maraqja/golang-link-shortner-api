// храним тут описание всех контрактов авторизации
package auth

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"` // добавляем теги для валидации
	Password string `json:"password"  validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
