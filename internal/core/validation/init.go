package core_validation

import (
	"github.com/go-playground/validator/v10"
)

var Validation *validator.Validate

type Validatable interface {
	Validate() error
}

func InitValidator[T ~string](roles []T, byteLength int) error {
	Validation = validator.New()

	err := Validation.RegisterValidation("pwd", func(fl validator.FieldLevel) bool {
		pwd := fl.Field().String()
		fn := fl.FieldName()

		if err := ValidatePassword(pwd, &fn); err != nil {
			return false
		}

		return true
	})
	if err != nil {
		return err
	}

	err = Validation.RegisterValidation("name", func(fl validator.FieldLevel) bool {
		name := fl.Field().String()
		fn := fl.FieldName()

		if err := ValidateStringLength(&name, fn, NameMinLen, NameMaxLen); err != nil {
			return false
		}

		return true
	})
	if err != nil {
		return err
	}

	err = Validation.RegisterValidation("user_role", func(fl validator.FieldLevel) bool {
		role := fl.Field().String()
		if err := ValidateArray(roles, T(role)); err != nil {
			return false
		}

		return true
	})
	if err != nil {
		return err
	}

	err = Validation.RegisterValidation("sid", func(fl validator.FieldLevel) bool {
		sid := fl.Field().String()
		if err := ValidateID(sid, byteLength); err != nil {
			return false
		}

		return true
	})
	if err != nil {
		return err
	}

	err = Validation.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		if _, err := ValidatePhoneNumber(phone); err != nil {
			return false
		}

		return true
	})
	if err != nil {
		return err
	}

	return nil
}
