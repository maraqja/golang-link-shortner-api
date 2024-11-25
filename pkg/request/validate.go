package request

import (
	"github.com/go-playground/validator/v10"
)

func Validate[T any](payload T) error {
	validator := validator.New()

	err := validator.Struct(payload) // валидация будет происходить по тегам
	return err
}
