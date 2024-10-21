package models

import (
	"github.com/go-playground/validator/v10"
)

type XValidator struct {
	*validator.Validate
}

var Validate = validator.New()

// InitValidator инициализирует валидатор.
func InitValidator() *XValidator {
	return &XValidator{
		Validate,
	}
}
