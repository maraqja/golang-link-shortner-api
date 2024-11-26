package main

import (
	"link-shortner-api/internal/link"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// db.Migrator().CreateTable()

	db.AutoMigrate(&link.Link{})
	/*
		db.AutoMigrate(&link.Link{}) автоматически создает или обновляет таблицу в базе данных на основе структуры Link.
		Вот что происходит:
			1. GORM анализирует структуру Link и все её поля
			2. Создает таблицу links (множественное число от Link)
			3. Добавляет все необходимые колонки:
				id, created_at, updated_at, deleted_at (из встроенной gorm.Model)
				url (из поля Url string)
				hash (из поля Hash string) с уникальным индексом (из-за тега gorm:"uniqueIndex")
			4. Если таблица уже существует, GORM проверит её структуру и добавит недостающие колонки/индексы
	*/
}
