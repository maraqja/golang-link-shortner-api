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

func TestLoginHandlerSuccess(t *testing.T) {
	handler, mock, err := bootstrap()
	if err != nil {
		t.Fatal(err)
	}
	rows := sqlmock.NewRows([]string{"email", "password", "name"}).
		AddRow("test@example.com", "$2a$10$GQS/HGBh4TyMG1.69MulZeL2B9uhOZWGfUow0PCD3.RkOigzYY0Dq", "kaka")
	// доавили виртуальную строку в моковую БД (пароль зашифрован)
	mock.ExpectQuery("SELECT (.+) FROM \"users\" WHERE email = ?").WillReturnRows(rows)

	data, err := json.Marshal(&auth.LoginRequest{
		Email:    "test@example.com",
		Password: "password",
	})
	if err != nil {
		t.Fatal(err)
	}

	reader := bytes.NewReader(data)
	wr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", reader)
	handler.Login()(wr, req)

	if wr.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, wr.Code)
	}
}

func TestRegisterHandlerSuccess(t *testing.T) {
	handler, mock, err := bootstrap()
	if err != nil {
		t.Fatal(err)
	}
	rows := sqlmock.NewRows([]string{"email", "password", "name"}) // для успешной регистрации не должен вернуть строк
	mock.ExpectQuery("SELECT").
		WillReturnRows(rows) // означает что на любой SELECT в рамках этого теста будут мокаться rows

	/*
		В GORM при выполнении операции Create (или других операций записи, таких как Update и Delete) автоматически открывается транзакция, если не отключили это поведение.
		Это делается для обеспечения целостности данных и возможности отката изменений в случае возникновения ошибки.
		Поэтому мокаем открытие и закрытие транзакции, тк нужно коммитнуть инсертнутую запись до вызова обработчика
	*/

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow((1)))
	mock.ExpectCommit()

	data, _ := json.Marshal(&auth.RegisterRequest{
		Email:    "new_test@example.com",
		Password: "new_test_password",
		Name:     "kaka",
	})

	reader := bytes.NewReader(data)
	wr := httptest.NewRecorder()                                          // создаем структуру записи ответа (респонса)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", reader) // создаем структуру запроса
	handler.Register()(wr, req)
	if wr.Code != http.StatusCreated {
		t.Fatalf("Expected status code %d, got %d", http.StatusCreated, wr.Code)
	}
}
