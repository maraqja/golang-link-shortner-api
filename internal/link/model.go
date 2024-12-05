package link

import (
	"link-shortner-api/internal/stat"
	"math/rand/v2"

	"gorm.io/gorm"
)

type Link struct {
	gorm.Model             // id, created_at, updated_at, deleted_at
	Url        string      `json:"url"`
	Hash       string      `json:"hash" gorm:"uniqueIndex"`                                       // hash должен быть уникален
	Stats      []stat.Stat `json:"stats" gorm:":constraint:OnUpdate:CASCADE, OnDelete:SET NULL;"` // связь с моделью Stat
	// при удалении ссылки у записи статистик, связанных с этой ссылкой link_id  будет равен NULL
	/*
		Почему не указываем явно связь с моделью Stat через foreignKey:LinkID?
		В GORM при определении отношений между моделями можно не указывать foreign key явно, потому что GORM использует соглашения об именовании (naming conventions).

		Когда мы определяем поле Stats []stat.Stat в модели Link, GORM автоматически:

		Определяет что это отношение один-ко-многим
		Создает foreign key с именем по шаблону [ИмяМодели]ID - в нашем случае LinkID
		Связывает поле LinkID в модели Stat с полем ID в модели Link
		Это работает потому что:

		В модели Stat есть поле LinkID uint
		Имя поля следует конвенции GORM
		Тип данных соответствует первичному ключу родительской модели (uint)
		Если бы мы хотели использовать другое имя для foreign key, тогда пришлось бы указать его явно через тег foreignKey.
	*/
}

func NewLink(url string) *Link {
	link := &Link{
		Url: url,
	}
	link.GenerateHash()
	return link
}

func (link *Link) GenerateHash() {
	link.Hash = randStringRunes(6)
}

// руна - это целое число (int32), которое представляет собой код символа в Unicode
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n) // создаем новый слайс рун
	for i := range b {
		b[i] = letterRunes[rand.IntN(len(letterRunes))] // заполняем результирующий слайс рандомно выбранными рунами из слайса letterRunes
	}
	return string(b) // возвращаем результат в виде строки
}
