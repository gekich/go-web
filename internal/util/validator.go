package util

import (
	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Validate(v any) error
}

type defaultValidator struct {
	v *validator.Validate
}

func New() Validator {
	return &defaultValidator{
		v: validator.New(),
	}
}

func (d *defaultValidator) Validate(v any) error {
	return d.v.Struct(v)
}
