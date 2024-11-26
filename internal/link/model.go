package link

import (
	"math/rand/v2"

	"gorm.io/gorm"
)

type Link struct {
	gorm.Model        // id, created_at, updated_at, deleted_at
	Url        string `json:"url"`
	Hash       string `json:"hash" gorm:"uniqueIndex"` // hash должен быть уникален
}

func NewLink(url string) *Link {
	return &Link{
		Url:  url,
		Hash: randStringRunes(6),
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n) // создаем новый слайс рун
	for i := range b {
		b[i] = letterRunes[rand.IntN(len(letterRunes))] // заполняем результирующий слайс рандомно выбранными рунами из слайса letterRunes
	}
	return string(b) // возвращаем результат в виде строки
}
