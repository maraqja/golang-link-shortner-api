package link

import (
	"link-shortner-api/pkg/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (repo *LinkRepository) GetById(id uint) (*Link, error) {
	var link Link
	result := repo.Database.DB.First(&link, id) // запишет результат в link
	if result.Error != nil {
		return nil, result.Error
	}
	return &link, nil
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

func (repo *LinkRepository) Update(link *Link) (*Link, error) {
	// сам понимает какую таблицу обновлять (тк это gorm model)
	result := repo.Database.DB.Clauses(clause.Returning{}).Updates(link)
	//		Clauses(clause.Returning{}) - нужен для того, чтобы вернуть обновленную запись (полная запись из postgres запишется в link)
	if result.Error != nil {
		return nil, result.Error
	}
	return link, nil
}

// удаляет запись (точнее делает soft-delete: устанавливает deleted_at)
// дальнейшие select-ы (при попытке перейти) не будут учитывать эту запись
func (repo *LinkRepository) Delete(id uint) error {

	// сам понимает какую таблицу обновлять (тк это gorm model)
	result := repo.Database.DB.Delete(&Link{}, id) // в delete указываем из какой таблицы удаляем + указать условие (по id если, то можно 2 аргументом просто id передать)
	// result := repo.Database.DB.Delete(&Link{}, "id = ?", id)
	// // Удаление по одному условию
	// repo.Database.DB.Where("hash = ?", hash).Delete(&Link{})
	// // Удаление по нескольким условиям
	// repo.Database.DB.Where("hash = ? AND active = ?", hash, true).Delete(&Link{})
	// // Через структуру условий
	// repo.Database.DB.Where(&Link{Hash: hash}).Delete(&Link{})
	if result.Error != nil {
		return result.Error
	}

	// RowsAffected вернет количество затронутых записей.
	// Если 0 - значит запись не была найдена.
	// Это более эффективный способ, так как используется один запрос к БД вместо двух.
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
