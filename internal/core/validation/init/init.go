package core_validation_init

import (
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
	"github.com/go-playground/validator/v10"
)

var Validation *validator.Validate

func InitValidator() error {
	Validation = validator.New()

	if err := Validation.RegisterValidation("pwd", registerPwdValidation); err != nil {
		return err
	}

	if err := Validation.RegisterValidation("name", registerNameValidation); err != nil {
		return err
	}

	if err := Validation.RegisterValidation("user_role", registerUserRoleValidation); err != nil {
		return err
	}

	if err := Validation.RegisterValidation("sort_type", registerSortTypeValidation); err != nil {
		return err
	}

	if err := Validation.RegisterValidation("sid", registerSessionIDValidation); err != nil {
		return err
	}

	if err := Validation.RegisterValidation("phone", registerPhoneValidation); err != nil {
		return err
	}

	return nil
}

func registerPwdValidation(fl validator.FieldLevel) bool {
	pwd := fl.Field().String()
	fn := fl.FieldName()

	if err := core_validation.ValidatePassword(pwd, &fn); err != nil {
		return false
	}

	return true
}

func registerNameValidation(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	fn := fl.FieldName()

	if err := core_validation.ValidateStringLength(
		&name, fn,
		core_validation.NameMinLen,
		core_validation.NameMaxLen,
	); err != nil {
		return false
	}

	return true
}

func registerUserRoleValidation(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	if err := core_validation.ValidateArray(domain.Roles, domain.UserRole(role)); err != nil {
		return false
	}

	return true
}

func registerSortTypeValidation(fl validator.FieldLevel) bool {
	sortType := fl.Field().String()
	if err := core_validation.ValidateArray(domain.SortTypes, domain.SortType(sortType)); err != nil {
		return false
	}

	return true
}

func registerSessionIDValidation(fl validator.FieldLevel) bool {
	sid := fl.Field().String()
	if err := core_validation.ValidateID(sid, domain.SessionIDByteLength); err != nil {
		return false
	}

	return true
}

func registerPhoneValidation(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if _, err := core_validation.ValidatePhoneNumber(phone); err != nil {
		return false
	}

	return true
}
