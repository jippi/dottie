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

func (d *Document) BelongsToGroup(name string) bool {
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

func (d *Document) GetGroup(name string) *Group {
	for _, grp := range d.Groups {
		if grp.BelongsToGroup(name) {
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

		group = doc.EnsureGroup(options.Group)

		existing = &Assignment{
			Name:    input.Name,
			Literal: input.Literal,
			Active:  input.Active,
			Group:   group,
		}

		if len(options.Before) > 0 {
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
		}

		if group != nil {
			group.Statements = append(group.Statements, existing)
		} else {
			idx := len(doc.Statements) - 1

			// if laste statement is a newline, replace it with the new assignment
			if idx > 1 && doc.Statements[idx].Is(&Newline{}) {
				doc.Statements[idx] = existing
			} else {
				// otherwise append it
				doc.Statements = append(doc.Statements, existing)
			}
		}
	}

	existing.Literal = input.Literal
	existing.Active = input.Active
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

	return true, nil
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

func (d *Document) GetPosition(name string) (int, *Assignment) {
	for i, assign := range d.Assignments() {
		if assign.Name == name {
			return i, assign
		}
	}

	return -1, nil
}
