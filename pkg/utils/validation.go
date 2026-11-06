package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"

	"jianco-api/pkg/apperror"
)

var validate = validator.New()

func Validate(value any) error {
	if err := validate.Struct(value); err != nil {
		var fieldErrors validator.ValidationErrors
		if ok := asValidationErrors(err, &fieldErrors); ok {
			parts := make([]string, 0, len(fieldErrors))
			for _, fieldError := range fieldErrors {
				parts = append(parts, fmt.Sprintf("%s failed %s validation", fieldError.Field(), fieldError.Tag()))
			}
			return apperror.BadRequest("validation failed: " + strings.Join(parts, ", "))
		}
		return apperror.BadRequest("validation failed")
	}
	return nil
}

func asValidationErrors(err error, target *validator.ValidationErrors) bool {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return false
	}
	*target = validationErrors
	return true
}
