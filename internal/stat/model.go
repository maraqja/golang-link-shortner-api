package stat

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Модель статистики имеет связь с ссылкой
// для каждой ссылки несколько статистик (количество кликовза разные дни типо)

type Stat struct { // Связь указываем в модели Link !!!
	gorm.Model
	LinkID uint           `json:"link_id" gorm:"not null"`
	Clicks int            `json:"clicks"`
	Date   datatypes.Date `json:"date"` // хотим хранить только дату (без времени)
}
