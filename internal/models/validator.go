package models

import (
	"github.com/go-playground/validator/v10"
)

type StructValidator struct {
	validate *validator.Validate
}

func NewValidator(v *validator.Validate) *StructValidator {
	return &StructValidator{
		validate: v,
	}
}

// Validator необходимо реализовать метод
func (v *StructValidator) Validate(out any) error {
	return v.validate.Struct(out)
}
