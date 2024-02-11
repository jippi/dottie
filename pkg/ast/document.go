// Package ast declares the types used to represent syntax trees for the .env file.
package ast

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/compose-spec/compose-go/template"
	"github.com/jippi/dottie/pkg/token"
)

// Document node represents .env file statement, that contains assignments and comments.
type Document struct {
	Statements  []Statement `json:"statements"` // Statements belonging to the root of the document
	Groups      []*Group    `json:"groups"`     // Groups within the document
	Annotations []*Comment  `json:"-"`          // Global annotations for configuration of dottie
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

func (d *Document) AllAssignments() []*Assignment {
	var assignments []*Assignment

	for _, statement := range d.Statements {
		if assign, ok := statement.(*Assignment); ok {
			assignments = append(assignments, assign)
		}
	}

	for _, group := range d.Groups {
		for _, statement := range group.Statements {
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

func (doc *Document) Interpolate(target *Assignment) (string, error) {
	if target == nil {
		return "", errors.New("can't interpolate a nil assignment")
	}

	if target.Quote.Is(token.SingleQuotes.Rune()) {
		return target.Literal, nil
	}

	lookup := func(input string) (string, bool) {
		// Lookup in process environment
		if val, ok := os.LookupEnv(input); ok {
			return val, ok
		}

		// Search the currently available assignments in the document
		result := doc.Get(input)
		if result == nil {
			return "", false
		}

		if !result.Active {
			return "", false
		}

		// If the assignment we found is on a line *after* the target
		// assignment, don't count it as found, since all normal shell interpolation
		// are handled in order (e.g. line 5 can't use a variable from line 10)
		if result.Position.Line >= target.Position.Line {
			return "", false
		}

		return result.Interpolated, true
	}

	return template.Substitute(target.Literal, lookup)
}

type UpsertPlacement uint

const (
	UpsertLast UpsertPlacement = iota
	UpsertAfter
	UpsertBefore
	UpsertFirst
)

type UpsertOptions struct {
	UpsertPlacementType  UpsertPlacement
	UpsertPlacementValue string
	Comments             []string
	ErrorIfMissing       bool
	Group                string
	SkipIfSame           bool
	SkipIfSet            bool
	SkipValidation       bool
}

func (doc *Document) Upsert(input *Assignment, options UpsertOptions) (*Assignment, error) {
	var group *Group

	existing := doc.Get(input.Name)

	if options.SkipIfSet && existing != nil && len(existing.Literal) > 0 && existing.Literal != "__CHANGE_ME__" && input.Literal != "__CHANGE_ME__" {
		return nil, nil
	}

	if options.SkipIfSame && existing != nil && existing.Literal == input.Literal && existing.Active == input.Active {
		return nil, nil
	}

	found := existing != nil

	// The key does not exists!
	if !found {
		if options.ErrorIfMissing {
			return nil, fmt.Errorf("Key [%s] does not exists", input.Name)
		}

		group = doc.EnsureGroup(options.Group)

		existing = &Assignment{
			Name:    input.Name,
			Literal: input.Literal,
			Active:  input.Active,
			Group:   group,
		}

		existingStatements := doc.Statements
		if existing.Group != nil {
			existingStatements = group.Statements
		}

		var res []Statement

		switch options.UpsertPlacementType {
		case UpsertFirst:
			res = append([]Statement{existing}, existingStatements...)

		case UpsertLast:
			res = append(existingStatements, existing)

		case UpsertAfter, UpsertBefore:
			for _, stmt := range existingStatements {
				assignment, ok := stmt.(*Assignment)
				if !ok {
					res = append(res, stmt)

					continue
				}

				switch {
				case options.UpsertPlacementType == UpsertBefore && assignment.Name == options.UpsertPlacementValue:
					res = append(res, existing, stmt)

				case options.UpsertPlacementType == UpsertAfter && assignment.Name == options.UpsertPlacementValue:
					res = append(res, stmt, existing)

				default:
					res = append(res, stmt)
				}
			}
		}

		if group != nil {
			group.Statements = res
		} else {
			doc.Statements = res
		}
	}

	if found {
		interpolated, err := doc.Interpolate(existing)
		if err != nil {
			return nil, errors.New("could not interpolate variable")
		}

		existing.Interpolated = interpolated
	}

	existing.Active = input.Active
	existing.Interpolated = input.Interpolated
	existing.Literal = input.Literal
	existing.Quote = input.Quote

	if comments := options.Comments; len(comments) > 0 {
		existing.Comments = nil

		for _, comment := range comments {
			if len(comment) == 0 && len(comments) == 1 {
				continue
			}

			existing.Comments = append(existing.Comments, NewComment(comment))
		}
	}

	if options.SkipValidation {
		return existing, nil
	}

	return existing, nil
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

func (d *Document) GetAssignmentIndex(name string) (int, *Assignment) {
	for i, assign := range d.Assignments() {
		if assign.Name == name {
			return i, assign
		}
	}

	return -1, nil
}
