package shared

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func validateNumeric(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	matched, _ := regexp.MatchString(`^[0-9]+$`, value)
	return matched
}

func RegisterCustomValidations(v *validator.Validate) {
	v.RegisterValidation("numeric", validateNumeric)
}
