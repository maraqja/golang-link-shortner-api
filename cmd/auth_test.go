package main

import (
	"bytes"
	"encoding/json"
	"io"
	"link-shortner-api/internal/auth"
	"link-shortner-api/internal/user"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// тесты выполняем в отдельной тестовой базе данных (предварительно создаем ее и выполняем миграции)
func initDb() *gorm.DB {
	err := godotenv.Load(".env") // читаем тестовый env
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// db.AutoMigrate(&link.Link{}, &user.User{}, &stat.Stat{})
	return db
}

func initData(db *gorm.DB) {
	db.Create(&user.User{
		Email:    "test@example.com",
		Password: "$2a$10$oYiAuyi1PddEFYOU8PhP2OI2Dvdnb0vYdsDV7Qgw4PtQPejMOkTeq", // зашифрованный пароль "password"
		Name: "kaka",
	})
}

func removeData(db *gorm.DB) {
	// db.Where("email = ?", "test@example.com").Delete(&user.User{}) // это не подойдет, тк будет soft-delete
	db.Where("email = ?", "test@example.com").Unscoped().Delete(&user.User{}) // удаляем полностью с помощью Unscoped
}

func TestLoginSuccess(t *testing.T) {
	db := initDb()
	initData(db)

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

	removeData(db)
}

func TestLoginFail(t *testing.T) {

	db := initDb()
	initData(db)
	
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

	removeData(db)
	
}