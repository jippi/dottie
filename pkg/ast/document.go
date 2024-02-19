// Package ast declares the types used to represent syntax trees for the .env file.
package ast

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jippi/dottie/pkg/template"
	"github.com/jippi/dottie/pkg/token"
	"go.uber.org/multierr"
)

// Document node represents .env file statement, that contains assignments and comments.
type Document struct {
	Statements  []Statement `json:"statements"` // Statements belonging to the root of the document
	Groups      []*Group    `json:"groups"`     // Groups within the document
	Annotations []*Comment  `json:"-"`          // Global annotations for configuration of dottie

	interpolateWarnings error
	interpolateErrors   error
}

func NewDocument() *Document {
	return &Document{}
}

func (d *Document) Is(other Statement) bool {
	if d == nil || other == nil {
		return false
	}

	return d.Type() == other.Type()
}

func (d *Document) Type() string {
	if d == nil {
		return "<nil>Document"
	}

	return reflect.TypeOf(d).String()
}

func (d *Document) BelongsToGroup(name string) bool {
	return false
}

func (d *Document) statementNode() {
}

func (d *Document) AllAssignments(selectors ...Selector) []*Assignment {
	var assignments []*Assignment

NEXT_STATEMENT:
	for _, statement := range d.Statements {
		// Filter
		for _, selector := range selectors {
			if selector(statement) == Exclude {
				continue NEXT_STATEMENT
			}
		}

		// Assign
		if assign, ok := statement.(*Assignment); ok {
			assignments = append(assignments, assign)
		}
	}

	for _, group := range d.Groups {
	NEXT_GROUP_STATEMENT:
		for _, statement := range group.Statements {
			// FILTER
			for _, selector := range selectors {
				if selector(statement) == Exclude {
					continue NEXT_GROUP_STATEMENT
				}
			}

			if assignment, ok := statement.(*Assignment); ok {
				assignments = append(assignments, assignment)
			}
		}
	}

	return assignments
}

func (d *Document) GetGroup(name string) *Group {
	for _, grp := range d.Groups {
		if grp.BelongsToGroup(name) {
			return grp
		}
	}

	return nil
}

func (d *Document) HasGroup(name string) bool {
	return d.GetGroup(name) != nil
}

func (d *Document) Get(name string) *Assignment {
	for _, assign := range d.AllAssignments() {
		if assign.Name == name {
			return assign
		}
	}

	return nil
}

func (d *Document) Has(name string) bool {
	return d.Get(name) != nil
}

func (doc *Document) InterpolateAll() (error, error) {
	defer func() {
		doc.interpolateWarnings = nil
		doc.interpolateErrors = nil
	}()

	for _, assignment := range doc.AllAssignments() {
		doc.doInterpolation(assignment)
	}

	return doc.interpolateWarnings, doc.interpolateErrors
}

func (doc *Document) InterpolateStatement(target *Assignment) (error, error) {
	defer func() {
		doc.interpolateWarnings = nil
		doc.interpolateErrors = nil
	}()

	doc.doInterpolation(target)

	return doc.interpolateWarnings, doc.interpolateErrors
}

func (doc *Document) doInterpolation(target *Assignment) {
	if target == nil {
		doc.interpolateErrors = multierr.Append(doc.interpolateErrors, errors.New("can't interpolate a nil assignment"))

		return
	}

	if !target.Enabled {
		return
	}

	target.Initialize()

	// Interpolate dependencies of the assignment before the assignment itself
	for _, dependency := range target.Dependencies {
		doc.doInterpolation(doc.Get(dependency.Name))
	}

	// If the assignment is wrapped in single quotes, no interpolation should happen
	if target.Quote.Is(token.SingleQuotes.Rune()) {
		target.Interpolated = target.Literal

		return
	}

	// If the assignment literal doesn't count any '$' it would never change from the
	// interpolated value
	if !strings.Contains(target.Literal, "$") {
		target.Interpolated = target.Literal

		return
	}

	value, warnings, err := template.Substitute(target.Literal, doc.interpolationMapper(target))
	if err != nil {
		err = fmt.Errorf("interpolation error for [%s] (%s): %w", target.Name, target.Position, err)
	}

	target.Interpolated = value

	doc.interpolateWarnings = multierr.Append(doc.interpolateWarnings, ContextualError(target, warnings))
	doc.interpolateErrors = multierr.Append(doc.interpolateErrors, ContextualError(target, err))
}

func (doc *Document) interpolationMapper(target *Assignment) func(input string) (string, bool) {
	return func(input string) (string, bool) {
		// Lookup in process environment
		if val, ok := os.LookupEnv(input); ok {
			return val, ok
		}

		// Search the currently available assignments in the document
		assignment := doc.Get(input)
		if assignment == nil {
			return "", false
		}

		// If the assignment we found is on a index (sorted) *after* the target
		// assignment, don't count it as found, since all normal shell interpolation
		// are handled in order (e.g. line 5 can't use a variable from line 10)
		if assignment.Position.Index >= target.Position.Index {
			return "", false
		}

		for _, dependency := range target.Dependencies {
			// Self-referencing is not allowed to avoid infinite loops in cases where you do [A="$A"]
			// which would trigger infinite recursive loop
			if dependency.Name == target.Name {
				doc.interpolateErrors = multierr.Append(doc.interpolateErrors, ContextualError(target, fmt.Errorf("Key [%s] must not reference itself", target.Name)))

				continue
			}

			// Lookup the assignment
			prerequisite := doc.Get(dependency.Name)

			// If it does not exists or is not enabled, abort
			if prerequisite == nil {
				doc.interpolateErrors = multierr.Append(doc.interpolateErrors, ContextualError(target, fmt.Errorf("Key [%s] must has invalid dependency [%s]", target.Name, dependency.Name)))

				continue
			}
		}

		return assignment.Interpolated, true
	}
}

func (doc *Document) EnsureGroup(name string) *Group {
	if len(name) == 0 {
		return nil
	}

	group := doc.GetGroup(name)

	if group == nil && len(name) > 0 {
		group = &Group{
			Name: "# " + name,
		}

		doc.Groups = append(doc.Groups, group)
	}

	return group
}

func (d *Document) GetConfig(name string) (string, error) {
	for _, comment := range d.Annotations {
		if comment.Annotation == nil {
			continue
		}

		if comment.Annotation.Key != name {
			continue
		}

		return comment.Annotation.Value, nil
	}

	return "", fmt.Errorf("could not find config key: [%s]", name)
}

func (d *Document) Assignments() []*Assignment {
	var assignments []*Assignment

	for _, statement := range d.Statements {
		if assign, ok := statement.(*Assignment); ok {
			assignments = append(assignments, assign)
		}
	}

	return assignments
}

func (d *Document) ReindexStatements() {
	for i, stmt := range d.AllAssignments() {
		stmt.Position.Index = i
	}
}

func (d *Document) GetAssignmentIndex(name string) (int, *Assignment) {
	for i, assign := range d.Assignments() {
		if assign.Name == name {
			return i, assign
		}
	}

	return -1, nil
}

func (document *Document) Initialize() {
	for _, assignment := range document.AllAssignments() {
		assignment.Initialize()

		// Add current assignment as dependent on its own dependencies
		for _, dependency := range assignment.Dependencies {
			if x := document.Get(dependency.Name); x != nil {
				if x.Dependents == nil {
					x.Dependents = make(map[string]*Assignment)
				}

				x.Dependents[assignment.Name] = assignment
			}
		}
	}

	document.ReindexStatements()
}

func (document *Document) Replace(assignment *Assignment) error {
	existing := document.Get(assignment.Name)
	if existing == nil {
		return fmt.Errorf("No KEY named [%s] exists in the document", assignment.Name)
	}

	if existing.Group != nil {
		for idx, stmt := range existing.Group.Statements {
			val, ok := stmt.(*Assignment)
			if !ok {
				continue
			}

			if val.Name == assignment.Name {
				existing.Group.Statements[idx] = assignment

				return nil
			}
		}
	}

	for idx, stmt := range document.Statements {
		val, ok := stmt.(*Assignment)
		if !ok {
			continue
		}

		if val.Name == assignment.Name {
			document.Statements[idx] = assignment

			return nil
		}
	}

	return fmt.Errorf("Could not find+replace KEY named [%s] in document", assignment.Name)
}

func (document *Document) Validate(selectors []Selector, ignoreErrors []string) ([]*ValidationError, error, error) {
	var (
		warnings, errors error
		data             = map[string]any{}
		rules            = map[string]any{}

		// The validation library uses a map[string]any as return value
		// which causes random ordering of keys. We would like them
		// to follow to order of which they are defined in the file
		// so this slice tracks that
		fieldOrder = []string{}
	)

NEXT:
	for _, assignment := range document.AllAssignments() {
		for _, selector := range selectors {
			status := selector(assignment)

			switch status {
			// Stop processing the statement and return nothing
			case Exclude:
				continue NEXT

			// Continue to next handler (or default behavior if we run out of handlers)
			case Keep:

			// Unknown signal
			default:
				panic(fmt.Errorf("unknown selector result: %v", status))
			}
		}

		validationRules := assignment.ValidationRules()
		if len(validationRules) == 0 {
			continue
		}

		warn, err := document.InterpolateStatement(assignment)

		warnings = multierr.Append(warnings, warn)
		errors = multierr.Append(errors, err)

		data[assignment.Name] = assignment.Interpolated
		rules[assignment.Name] = validationRules

		fieldOrder = append(fieldOrder, assignment.Name)
	}

	validationErrors := validator.New().ValidateMap(data, rules)

	var result []*ValidationError

NEXT_FIELD:
	for _, field := range fieldOrder {
		err, ok := validationErrors[field]
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

		result = append(result, &ValidationError{
			WrappedError: err,
			Assignment:   document.Get(field),
		})
	}

	return result, warnings, errors
}

func (document *Document) ValidateSingleAssignment(assignment *Assignment, selectors []Selector, ignoreErrors []string) (ValidationErrors, error, error) {
	return document.Validate(
		append(
			[]Selector{
				ExcludeDisabledAssignments,
				RetainExactKey(assignment.Name),
			},
			selectors...,
		),
		ignoreErrors,
	)
}
