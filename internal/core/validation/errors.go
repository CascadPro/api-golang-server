package core_validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func MapErrors(err error) []error {
	errs := make([]error, 0)

	verrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return append(errs, err)
	}

	for _, fe := range verrs {
		switch fe.Tag() {
		case "required":
			errs = append(errs, fmt.Errorf("`%s` can't be NULL", fe.Field()))
		case "email":
			errs = append(errs, fmt.Errorf("`%s` must be valid email address", fe.Field()))
		case "phone":
			errs = append(errs, fmt.Errorf("`%s` must be valid phone number", fe.Field()))
		case "name":
			errs = append(errs, fmt.Errorf("`%s` must be valid '%s' string", fe.Field(), fe.Tag()))
		case "min":
			errs = append(errs, fmt.Errorf("`%s` minimum length is %s symbols / items", fe.Field(), fe.Param()))
		case "max":
			errs = append(errs, fmt.Errorf("`%s` maximum length is %s symbols / items", fe.Field(), fe.Param()))
		case "unique":
			errs = append(errs, fmt.Errorf("`%s` must contain only unique items", fe.Field()))
		case "number":
			errs = append(errs, fmt.Errorf("`%s` must be valid number", fe.Field()))
		case "lte":
			errs = append(errs, fmt.Errorf("`%s` must be less than or equal %s", fe.Field(), fe.Param()))
		case "gte":
			errs = append(errs, fmt.Errorf("`%s` must be greater than or equal %s", fe.Field(), fe.Param()))
		case "uuid":
			errs = append(errs, fmt.Errorf("`%s` must be valid universal user identifier", fe.Field()))
		case "pwd":
			errs = append(errs, fmt.Errorf("`%s` must be complex password string", fe.Field()))
		case "user_role":
			errs = append(errs, fmt.Errorf("Field `%s` with value '%s' isn't valid `UserRole`", fe.Field(), fe.Value()))
		case "sort_type":
			errs = append(errs, fmt.Errorf("Field `%s` with value '%s' isn't valid `SortType`", fe.Field(), fe.Value()))
		default:
			errs = append(errs, fmt.Errorf("`%s failed validation, tag: %s", fe.Field(), fe.Tag()))
		}
	}

	return errs
}
