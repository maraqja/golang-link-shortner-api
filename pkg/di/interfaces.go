package di

import "link-shortner-api/internal/user"

// Интерфейсы, которые могут бы использованы для внедрения зависимостей
type IStatRepository interface { // не зависим от конкретной имплементации
	AddClick(linkId uint)
}

type IUserRepository interface {
	Create(user *user.User) (*user.User, error)
	FindByEmail(email string) (*user.User, error)
}
