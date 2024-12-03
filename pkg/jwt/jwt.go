package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWT struct { // по сути это скорее конфиг для работы с JWT
	Secret string
}

func NewJWT(secret string) *JWT { // а это конструктор конфига
	return &JWT{
		Secret: secret,
	}
}

func (j *JWT) Create(email string) (string, error) {
	// под клеймами понимаются свойства в payload-объекте
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
	}) // сгенерили токен (структуру(объект) токена, без указанного свойства sign)
	signed_token, err := token.SignedString([]byte(j.Secret)) // формируем сам jwt (как строку)
	if err != nil {
		return "", err
	}
	return signed_token, nil
}
