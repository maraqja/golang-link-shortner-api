package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTData struct {
	Email string
}

type JWT struct { // по сути это скорее конфиг для работы с JWT
	Secret string
}

func NewJWT(secret string) *JWT { // а это конструктор конфига
	return &JWT{
		Secret: secret,
	}
}

func (j *JWT) Create(data JWTData) (string, error) {
	// под клеймами понимаются свойства в payload-объекте
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": data.Email,
	}) // сгенерили токен (структуру(объект) токена, без указанного свойства sign)
	signed_token, err := token.SignedString([]byte(j.Secret)) // формируем сам jwt (как строку)
	if err != nil {
		return "", err
	}
	return signed_token, nil
}

func (j *JWT) Parse(token string) (bool, *JWTData) { //
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})
	if err != nil {
		return false, nil
	}

	// MapClaims - map[string]interface{}
	email := t.Claims.(jwt.MapClaims)["email"] // п
	return t.Valid, &JWTData{
		Email: email.(string),
	}
}
