package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidationError wraps multiple validation errors
type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	return strings.Join(e.Errors, "; ")
}

func ValidateStruct(data any) error {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	var validationErrors []string
	for _, e := range err.(validator.ValidationErrors) {
		fieldName := e.Field()
		if jsonTag := e.StructField(); jsonTag != "" {
			fieldName = strings.ToLower(e.Field())
		}

		switch e.Tag() {
		case "required":
			validationErrors = append(validationErrors, fmt.Sprintf("%s is required", fieldName))
		case "email":
			validationErrors = append(validationErrors, fmt.Sprintf("%s must be a valid email", fieldName))
		case "min":
			validationErrors = append(validationErrors, fmt.Sprintf("%s must be at least %s characters", fieldName, e.Param()))
		case "max":
			validationErrors = append(validationErrors, fmt.Sprintf("%s must be at most %s characters", fieldName, e.Param()))
		case "oneof":
			validationErrors = append(validationErrors, fmt.Sprintf("%s must be one of: %s", fieldName, e.Param()))
		case "url":
			validationErrors = append(validationErrors, fmt.Sprintf("%s must be a valid URL", fieldName))
		default:
			validationErrors = append(validationErrors, fmt.Sprintf("%s is invalid", fieldName))
		}
	}

	if len(validationErrors) > 0 {
		return &ValidationError{Errors: validationErrors}
	}

	return errors.New("validation failed")
}
