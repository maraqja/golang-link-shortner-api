package jwt_test

import (
	"link-shortner-api/pkg/jwt"
	"testing"
)

func TestJWTCreate(t *testing.T) {
	jwtService := jwt.NewJWT("1ERXjnPD2YViiEFTMevuHafu0SJFnsEK5jFJgvb-NvA")
	const email = "test@example.com"
	token, err := jwtService.Create(jwt.JWTData{
		Email: email,
	})

	if err != nil {
		t.Fatal(err)
	}

	isValid, data := jwtService.Parse(token)
	if !isValid {
		t.Fatal("Expected valid token")
	}

	if data.Email != email {
		t.Fatalf("Expected email %s, got %s", email, data.Email)
	}



}