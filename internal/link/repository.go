package link

import (
	"link-shortner-api/pkg/db"
)

// тут реализация взаимодействия с БД
// инжектируем репозиторий в сервисы
// паттерн репозитория лучше с архитектурной точки зрения + позволяет в любой момент заменить БД

type LinkRepository struct {
	Database *db.Db
}

func NewLinkRepository(database *db.Db) *LinkRepository {
	return &LinkRepository{
		Database: database,
	}
}

func (repo *LinkRepository) Create(link *Link) error {

}
