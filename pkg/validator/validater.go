package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var V *validator.Validate

func Init() {
	V = validator.New()
}

// FormatErrors converts validator errors into a readable message
func FormatErrors(err error) string {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return err.Error()
	}

	var msgs []string
	for _, fe := range ve {
		msgs = append(msgs, formatFieldError(fe))
	}
	return strings.Join(msgs, "; ")
}

func formatFieldError(fe validator.FieldError) string {
	field := strings.ToLower(fe.Field()[:1]) + fe.Field()[1:]

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, fe.Param())
	case "gte":
		return fmt.Sprintf("%s must be >= %s", field, fe.Param())
	case "lte":
		return fmt.Sprintf("%s must be <= %s", field, fe.Param())
	case "gt":
		return fmt.Sprintf("%s must be > %s", field, fe.Param())
	case "lt":
		return fmt.Sprintf("%s must be < %s", field, fe.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of [%s]", field, fe.Param())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}