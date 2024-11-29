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
