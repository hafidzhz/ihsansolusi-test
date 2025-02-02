package validator

import (
	"fmt"

	vldtr "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func NewValidator() *vldtr.Validate {
	validate := vldtr.New()

	_ = validate.RegisterValidation("uuid", func(fl vldtr.FieldLevel) bool {
		field := fl.Field().String()
		if _, err := uuid.Parse(field); err != nil {
			return true
		}
		return false
	})

	return validate
}

func ValidatorErrors(err error) map[string]string {
	fields := map[string]string{}

	for _, err := range err.(vldtr.ValidationErrors) {
		errMsg := fmt.Sprintf("validation failed on '%s' tag", err.Tag())
		param := err.Param()
		if param != "" {
			errMsg = fmt.Sprintf("%s. allow: %s", errMsg, param)
		}
		fields[err.Field()] = errMsg
	}

	return fields
}
