package config

import "github.com/go-playground/validator/v10"

var Validate *validator.Validate

//this will instantiate validator package globally
func InstantiateValidator() {
	Validate = validator.New()
}
