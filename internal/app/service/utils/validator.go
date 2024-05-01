package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func hasUppercase(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`[A-Z]`).MatchString(fl.Field().String())
}

func hasLowercase(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`[a-z]`).MatchString(fl.Field().String())
}

func hasNumber(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`[0-9]`).MatchString(fl.Field().String())
}

func hasSpecialChar(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(fl.Field().String())
}

var Validate *validator.Validate

func InitValidator() {
	Validate = validator.New()

	Validate.RegisterValidation("uppercase", hasUppercase)
	Validate.RegisterValidation("lowercase", hasLowercase)
	Validate.RegisterValidation("number", hasNumber)
	Validate.RegisterValidation("specialchar", hasSpecialChar)
}