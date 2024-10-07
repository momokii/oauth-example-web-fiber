package utils

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

// GetValidator returns a singleton instance of the validator
func GetValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New(validator.WithRequiredStructEnabled())
	})
	return validate
}

// ValidateStruct validates a struct
func ValidateStruct(s interface{}) error {
	return GetValidator().Struct(s)
}
