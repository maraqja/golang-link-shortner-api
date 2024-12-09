package auth_test

import (
	"link-shortner-api/internal/auth"
	"link-shortner-api/internal/user"
	"testing"
)

// Интеграционные тесты (не e2e, тк БД не поднимаем например, а просто мокаем репозиторий)
type MockUserRepository struct{}

func (repo *MockUserRepository) Create(u *user.User) (*user.User, error) {
	return &user.User{
		Email: "test@mail.ru",
	}, nil
} // для успешной регистрации должны создать пользователя без ошибки
func (repo *MockUserRepository) FindByEmail(email string) (*user.User, error) {
	return nil, nil
} // для успешной регистрации не должны найти пользователя

func TestRegisterSuccess(t *testing.T) {
	const initialEmail = "test@mail.ru"
	authService := auth.NewAuthService(&MockUserRepository{}) // DI через интерфейсы позволяет завязываться на необходимость наличия методов определенных, а не на конкретную реализацию (иначе бы пришлось создать репозиторий, БД для репозитория чтобы внедрить в качестве зависимости в сервис)
	email, err := authService.Register(initialEmail, "password", "kaka")
	if err != nil {
		t.Fatal(err)
		return
	}
	if email != initialEmail {
		t.Fatalf("Email %s do not match %s", email, initialEmail)
		return
	}
}
