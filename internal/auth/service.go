package auth

import (
	"errors"
	"link-shortner-api/internal/user"
)

type AuthService struct {
	UserRepository *user.UserRepository
}

func NewAuthService(userRepository *user.UserRepository) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
	}
}

func (service *AuthService) Register(email, password, name string) (string, error) { // не используем входной парметр типа auth.RegisterRequest тк у сервиса и хендлера (контроллера) должны быть разные контракты (можно потом вынести в отдельный интерфейс)
	existedUser, _ := service.UserRepository.FindByEmail(email)
	if existedUser != nil {
		return "", errors.New(ErrUserExists)
	}
	user := &user.User{
		Email:    email,
		Password: "", // надо захэшировать
		Name:     name,
	}
	_, err := service.UserRepository.Create(user)
	/*
		Почему лучше не сначала создать структуру и передавать указатель только в функцию
		// 1. Сначала создается структура как значение
		user := user.User{
			Email: email,
			Name: name,
		}
		// Здесь структура создается в стеке

		// 2. Затем при передаче &user создается новый указатель на эту структуру
		service.UserRepository.Create(&user)
	*/
	if err != nil {
		return "", err
	}
	return user.Email, nil
}
