package core_validation

import (
	"github.com/go-playground/validator/v10"
)

var Validation *validator.Validate

type Validatable interface {
	Validate() error
}

func InitValidator[T ~string](roles []T) error {
	Validation = validator.New()

	return nil
}
