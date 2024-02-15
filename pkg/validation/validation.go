package validation

import (
	"context"
	"fmt"
	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/render"
)

type ValidationError struct {
	WrappedError any
	Assignment   *ast.Assignment
}

func (e ValidationError) Error() string {
	if val, ok := e.WrappedError.(error); ok {
		return val.Error()
	}

	return fmt.Sprintf("%+v", e.WrappedError)
}

func NewError(assignment *ast.Assignment, err error) ValidationError {
	return ValidationError{
		WrappedError: err,
		Assignment:   assignment,
	}
}

func Validate(ctx context.Context, doc *ast.Document, handlers []render.Handler, ignoreErrors []string) []ValidationError {
	data := map[string]any{}
	rules := map[string]any{}

	// The validation library uses a map[string]any as return value
	// which causes random ordering of keys. We would like them
	// to follow to order of which they are defined in the file
	// so this slice tracks that
	fieldOrder := []string{}

NEXT:
	for _, assignment := range doc.AllAssignments() {
		handlerInput := &render.HandlerInput{
			CurrentStatement:  assignment,
			PreviousStatement: nil,
			Renderer:          nil,
			Settings:          *render.NewSettings(),
		}

		for _, handler := range handlers {
			status := handler(ctx, handlerInput)

			switch status {
			// Stop processing the statement and return nothing
			case render.Stop:
				continue NEXT

			// Continue to next handler (or default behavior if we run out of handlers)
			case render.Continue, render.Return:

			// Unknown signal
			default:
				panic("unknown signal: " + status.String())
			}
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

NEXT_FIELD:
	for _, field := range fieldOrder {
		err, ok := errors[field]
		if !ok {
			continue
		}

		switch err := err.(type) {
		case validator.ValidationErrors:
			for _, rule := range err {
				if slices.Contains(ignoreErrors, rule.ActualTag()) {
					continue NEXT_FIELD
				}
			}
		}

		result = append(result, ValidationError{
			WrappedError: err,
			Assignment:   doc.Get(field),
		})
	}

	return result
}

func ValidateSingleAssignment(ctx context.Context, doc *ast.Document, assignment *ast.Assignment, handlers []render.Handler, ignoreErrors []string) []ValidationError {
	keys := AssignmentsToValidateRecursive(assignment)

	return Validate(
		ctx,
		doc,
		append(
			[]render.Handler{
				render.ExcludeDisabledAssignments,
				render.RetainExactKey(keys...),
			},
			handlers...,
		),
		ignoreErrors,
	)
}

func AssignmentsToValidateRecursive(assignment *ast.Assignment) []string {
	keys := []string{assignment.Name}

	for _, dependent := range assignment.Dependents {
		keys = append(keys, dependent.Name)

		if len(dependent.Dependents) > 0 {
			for _, d := range dependent.Dependents {
				keys = append(keys, AssignmentsToValidateRecursive(d)...)
			}
		}
	}

	return slices.Compact(keys)
}
