package handlers

import "github.com/go-playground/validator"

var validateStruct *validator.Validate

func init() {
	validateStruct = validator.New()
}
