package main

import (
	"bytes"
	"encoding/json"
	"io"
	"link-shortner-api/internal/auth"
	"net/http"
	"net/http/httptest"
	"testing"
)



func TestLoginSuccess(t *testing.T) {
	testServer := httptest.NewServer(App())
	defer testServer.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "test@example.com",
		Password: "password",
	})
	res, err := http.Post(testServer.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	var resData auth.LoginResponse
	err = json.Unmarshal(body, &resData)
	if err != nil {
		t.Fatal(err)
	}
	if resData.Token == "" {
		t.Fatal("Expected token")
	}
}

func TestLoginFail(t *testing.T) {
	testServer := httptest.NewServer(App())
	defer testServer.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "test@example.com",
		Password: "invalid_password",
	})
	res, err := http.Post(testServer.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected status code %d, got %d", http.StatusUnauthorized, res.StatusCode)
	}
	
}