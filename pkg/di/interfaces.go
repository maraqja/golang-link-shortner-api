package di

// Интерфейсы, которые могут бы использованы для внедрения зависимостей
type IStatRepository interface { // не зависим от конкретной имплементации
	AddClick(linkId uint)
}
