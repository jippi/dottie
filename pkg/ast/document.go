// Package ast declares the types used to represent syntax trees for the .env file.
package ast

import (
	"fmt"
	"os"
	"reflect"
)

// Document node represents .env file statement, that contains assignments and comments.
type Document struct {
	Statements  []Statement `json:"statements"` // Statements belonging to the root of the document
	Groups      []*Group    `json:"groups"`     // Groups within the document
	Annotations []*Comment  `json:"-"`          // Global annotations for configuration of dottie
}

func (d *Document) Is(other Statement) bool {
	return reflect.TypeOf(d) == reflect.TypeOf(other)
}

func (d *Document) BelongsToGroup(config RenderSettings) bool {
	return false
}

func (d *Document) statementNode() {
}

func (d *Document) Assignments() []*Assignment {
	var values []*Assignment

	for _, stmt := range d.Statements {
		if assign, ok := stmt.(*Assignment); ok {
			values = append(values, assign)
		}
	}

	for _, grp := range d.Groups {
		for _, stmt := range grp.Statements {
			if assign, ok := stmt.(*Assignment); ok {
				values = append(values, assign)
			}
		}
	}

	return values
}

func (d *Document) GetGroup(config RenderSettings) *Group {
	for _, grp := range d.Groups {
		if grp.BelongsToGroup(config) {
			return grp
		}
	}

	return nil
}

func (d *Document) Get(name string) *Assignment {
	for _, assign := range d.Assignments() {
		if assign.Name == name {
			return assign
		}
	}

	return nil
}

func (d *Document) GetInterpolation(in string) (string, bool) {
	// Lookup in process environment
	if val, ok := os.LookupEnv(in); ok {
		return val, ok
	}

	// Search the currently available assignments in the document
	assignment := d.Get(in)
	if assignment == nil {
		return "", false
	}

	if !assignment.Active {
		return "", false
	}

	return assignment.Interpolated, true
}

type SetOptions struct {
	ErrorIfMissing bool
	SkipIfSet      bool
	SkipIfSame     bool
	Group          string
	Before         string
	Comments       []string
}

func (doc *Document) Set(input *Assignment, options SetOptions) (bool, error) {
	var group *Group

	existing := doc.Get(input.Name)

	if options.SkipIfSet && existing != nil && len(existing.Literal) > 0 && existing.Literal != "__CHANGE_ME__" && input.Literal != "__CHANGE_ME__" {
		return false, nil
	}

	if options.SkipIfSame && existing != nil && len(existing.Literal) > 0 && existing.Literal == input.Literal {
		return false, nil
	}

	// The key does not exists!
	if existing == nil {
		if options.ErrorIfMissing {
			return false, fmt.Errorf("Key [%s] does not exists", input.Name)
		}

		group = doc.GetGroup(RenderSettings{FilterGroup: options.Group})
		if group == nil {
			group = &Group{
				Name: input.Group.Name,
			}

			doc.Groups = append(doc.Groups, group)
		}

		existing = &Assignment{
			Name:  input.Name,
			Group: group,
		}

		switch {
		case len(options.Before) > 0:
			before := options.Before

			var res []Statement
			for _, stmt := range group.Statements {
				x, ok := stmt.(*Assignment)
				if !ok {
					res = append(res, stmt)
					continue
				}

				if x.Name == before {
					res = append(res, existing)
				}

				res = append(res, stmt)
			}

			group.Statements = res

		default:
			group.Statements = append(group.Statements, existing)
		}
	}

	existing.Literal = input.Literal
	existing.Active = input.Active
	existing.Quote = input.Quote

	if comments := options.Comments; len(comments) > 0 {
		existing.Comments = nil

		for _, comment := range comments {
			existing.Comments = append(existing.Comments, NewComment(comment))
		}
	}

	return true, nil
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

func (d *Document) GetPosition(name string) (int, *Assignment) {
	for i, assign := range d.Assignments() {
		if assign.Name == name {
			return i, assign
		}
	}

	return -1, nil
}

func (d *Document) RenderFull() string {
	return d.Render(RenderSettings{
		IncludeCommented: true,
		ShowBlankLines:   true,
		ShowColors:       false,
		ShowComments:     true,
		ShowGroups:       true,
		Interpolate:      false,
	})
}

func (d *Document) Render(config RenderSettings) string {
	return renderStatements(d.Statements, config)
}
