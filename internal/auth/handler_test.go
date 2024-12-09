package auth_test

import (
	"link-shortner-api/configs"
	"link-shortner-api/internal/auth"
	"link-shortner-api/internal/user"
	"link-shortner-api/pkg/db"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestLoginSuccess(t *testing.T) {
	// go get github.com/DATA-DOG/go-sqlmock
	database, mock, err := sqlmock.New() // создали мок БД
	if err != nil {
		t.Fatalf("Failed init mock db: %s: ", err)
		return
	}
	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: database,
	})) // получили полноценный инстанс gorm, но при этом не поднимаем реальную БД (мокнули БД)
	if err != nil {
		t.Fatalf("Failed init gorm: %s: ", err)
		return
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

}
