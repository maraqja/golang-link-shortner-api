package auth_test

import (
	"bytes"
	"encoding/json"
	"link-shortner-api/configs"
	"link-shortner-api/internal/auth"
	"link-shortner-api/internal/user"
	"link-shortner-api/pkg/db"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func bootstrap() (*auth.AuthHandler, sqlmock.Sqlmock, error) {
	database, mock, err := sqlmock.New() // создали мок БД
	if err != nil {
		return nil, nil, err
	}
	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: database,
	})) // получили полноценный инстанс gorm, но при этом не поднимаем реальную БД (мокнули БД)
	if err != nil {
		return nil, nil, err
	}
	defer database.Close()
	userRepo := user.NewUserRepository(&db.Db{ // создали репозиторий с мокнутой БД
		DB: gormDb,
	})

	handler := auth.AuthHandler{
		Config: &configs.Config{
			Auth: configs.AuthConfig{
				Secret: "secret",
			},
		},
		AuthService: auth.NewAuthService(userRepo),
	}
	return &handler, mock, nil
}

func TestLoginSuccess(t *testing.T) {
	handler, mock, err := bootstrap()
	if err != nil {
		t.Fatal(err)
	}
	rows := sqlmock.NewRows([]string{"email", "password"}).
		AddRow("test@example.com", "$2a$10$GQS/HGBh4TyMG1.69MulZeL2B9uhOZWGfUow0PCD3.RkOigzYY0Dq")
	// доавили виртуальную строку в моковую БД (пароль зашифрован)

	mock.ExpectQuery("SELECT email, password FROM \"users\" WHERE email = ?").
		WithArgs("test@example.com").
		WillReturnRows(rows) // указываем при каком запросе (в данном случае выборка по email) какие данные должны быть возвращены
	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "test@example.com",
		Password: "password",
	})

	reader := bytes.NewReader(data)
	wr := httptest.NewRecorder()                                       // создаем структуру записи ответа (респонса)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", reader) // создаем структуру запроса
	handler.Login()(wr, req)

	if wr.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, wr.Code)
	}
}
