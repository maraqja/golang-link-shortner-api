package link

import (
	"fmt"
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
	fmt.Printf("link = %v", link)
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

func (repo *LinkRepository) GetCount() (int64, error) {
	var count int64
	result := repo.Database.DB.
		Table("links").
		Where("deleted_at IS NULL").
		Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (repo *LinkRepository) GetAll(limit, offset int) ([]Link, error) {
	var links []Link

	// result := repo.Database.DB.Limit(limit).Offset(offset).Find(&links)
	// Используем query builder
	result := repo.Database.DB.
		Table("links").
		Select("*"). // если указываем не все поля, то нужно описать новую структуру
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(&links) // разница с Find: Scan не заполняет gorm.Model

	/*
			// Разница Find и Scan:



			type Link struct {
				gorm.Model   // (ID, CreatedAt, UpdatedAt, DeletedAt)
				Hash string
				URL  string
			}

			// Используем Find
			var link1 Link
			db.Table("links").Where("id = 1").Find(&link1)
			// После этого если сделаем:
			link1.Hash = "new-hash"
			db.Save(&link1)  // GORM автоматически обновит UpdatedAt

			// Используем Scan
			var link2 Link
			db.Table("links").Where("id = 1").Scan(&link2)
			// После этого если сделаем:
			link2.Hash = "new-hash"
			db.Save(&link2)  // UpdatedAt придется обновлять вручную

			// То есть Find() подготавливает структуру для дальнейшей работы с GORM (делает её "отслеживаемой"),
			 а Scan() просто копирует данные из БД в структуру без дополнительной подготовки.

		/*

		/*
					// Если нужно использоваться фулл запрос
				result := repo.Database.DB.Raw(`
			        SELECT *
			        FROM links
			        WHERE deleted_at IS NULL
			        ORDER BY created_at DESC
			        LIMIT ?
			        OFFSET ?
			    `, limit, offset).Scan(&links)
	*/

	if result.Error != nil {
		return nil, result.Error

	}
	return links, nil
}
