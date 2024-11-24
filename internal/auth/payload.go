// храним тут описание всех контрактов авторизации
package auth

type LoginResponse struct {
	Token string `json:"token"`
}