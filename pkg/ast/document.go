// Package ast declares the types used to represent syntax trees for the .env file.
package ast

import (
	"fmt"
	"reflect"
)

// Document node represents .env file statement, that contains assignments and comments.
type Document struct {
	Statements  []Statement   `json:"statements"`
	Groups      []*Group      `json:"groups"`
	Assignments []*Assignment `json:"-"`
	Comments    []*Comment    `json:"-"`
}

func (d *Document) Is(other Statement) bool {
	return reflect.TypeOf(d) == reflect.TypeOf(other)
}

func (d *Document) BelongsToGroup(config RenderSettings) bool {
	return false
}

func (d *Document) statementNode() {
}

func (d *Document) AllAssignments() []*Assignment {
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
	for _, assign := range d.Assignments {
		if assign.Key == name {
			return assign
		}
	}

	return nil
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

	existing := doc.Get(input.Key)

	if options.SkipIfSet && existing != nil && len(existing.Value) > 0 && existing.Value != "__CHANGE_ME__" && input.Value != "__CHANGE_ME__" {
		return false, nil
	}

	if options.SkipIfSame && existing != nil && len(existing.Value) > 0 && existing.Value == input.Value {
		return false, nil
	}

	// The key does not exists!
	if existing == nil {
		if options.ErrorIfMissing {
			return false, fmt.Errorf("Key [%s] does not exists", input.Key)
		}

		group = doc.GetGroup(RenderSettings{FilterGroup: options.Group})
		if group == nil {
			group = &Group{
				Name: input.Group.Name,
			}

			doc.Groups = append(doc.Groups, group)
		}

		existing = &Assignment{
			Key:   input.Key,
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

				if x.Key == before {
					res = append(res, existing)
				}

				res = append(res, stmt)
			}

			group.Statements = res

		default:
			group.Statements = append(group.Statements, existing)
		}
	}

	doc.Assignments = append(doc.Assignments, existing)

	existing.Value = input.Value
	existing.Commented = input.Commented
	existing.Quoted = input.Quoted

	if comments := options.Comments; len(comments) > 0 {
		existing.Comments = nil

		for _, comment := range comments {
			existing.Comments = append(existing.Comments, NewComment(comment))
		}
	}

	return true, nil
}

func (d *Document) GetConfig(name string) (string, error) {
	for _, comment := range d.Comments {
		if !comment.Annotation {
			continue
		}

		if comment.AnnotationKey != name {
			continue
		}

		return comment.AnnotationValue, nil
	}

	return "", fmt.Errorf("could not find config key: [%s]", name)
}

func (d *Document) GetPosition(name string) (int, *Assignment) {
	for i, assign := range d.Assignments {
		if assign.Key == name {
			return i, assign
		}
	}

	return -1, nil
}

func (d *Document) RenderFull() string {
	return d.Render(RenderSettings{
		ShowPretty:       true,
		IncludeCommented: true,
	})
}

func (d *Document) Render(config RenderSettings) string {
	return renderStatements(d.Statements, config)
}
