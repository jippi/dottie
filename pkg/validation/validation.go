package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/jippi/dottie/pkg/ast"
)

type ValidationError struct {
	Error      any
	Assignment *ast.Assignment
}

func NewError(assignment *ast.Assignment, err error) ValidationError {
	return ValidationError{
		Error:      err,
		Assignment: assignment,
	}
}

func Validate(d *ast.Document) []ValidationError {
	data := map[string]any{}
	rules := map[string]any{}

	// The validation library uses a map[string]any as return value
	// which causes random ordering of keys. We would like them
	// to follow to order of which they are defined in the file
	// so this slice tracks that
	fieldOrder := []string{}

	for _, assignment := range d.Assignments() {
		if !assignment.Active {
			continue
		}

		validationRules := assignment.ValidationRules()
		if len(validationRules) == 0 {
			continue
		}

		data[assignment.Name] = assignment.Interpolated
		rules[assignment.Name] = validationRules

		fieldOrder = append(fieldOrder, assignment.Name)
	}

	errors := validator.New(validator.WithRequiredStructEnabled()).ValidateMap(data, rules)

	result := []ValidationError{}

	for _, field := range fieldOrder {
		if err, ok := errors[field]; ok {
			result = append(result, ValidationError{
				Error:      err,
				Assignment: d.Get(field),
			})
		}
	}

	return result
}
