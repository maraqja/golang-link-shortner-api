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

func (repo *LinkRepository) Create(link *Link) (*Link, error) {
	result := repo.Database.DB.Create(link) // для создания не нужно указывать таблицу, GORM сам понимает куда вставлять по переданной структуре, тк Link содержит поле gorm.Model
	// Create обогатит link данными gorm.Model (тк передаем по указателю)
	if result.Error != nil {
		return nil, result.Error
	}
	return link, nil
}

func (repo *LinkRepository) GetByHash(hash string) (*Link, error) {
	var link Link

	// repo.Database.DB.First(&link, "hash = ? OR id = ?", hash, id)
	result := repo.Database.DB.Where("hash = ?", hash).First(&link) // запишет результат в link
	if result.Error != nil {
		return nil, result.Error
	}
	return &link, nil
}
