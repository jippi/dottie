// Package ast declares the types used to represent syntax trees for the .env file.
package ast

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

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

	// If the assignment is wrapped in single quotes, no interpolation should happen
	if target.Quote.Is(token.SingleQuotes.Rune()) {
		return target.Literal, nil
	}

	// If the assignment literal doesn't count any '$' it would never change from the
	// interpolated value
	if !strings.Contains(target.Literal, "$") {
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

		if !result.Enabled {
			return "", false
		}

		// If the assignment we found is on a index (sorted) *after* the target
		// assignment, don't count it as found, since all normal shell interpolation
		// are handled in order (e.g. line 5 can't use a variable from line 10)
		if result.Position.Index >= target.Position.Index {
			return "", false
		}

		return result.Interpolated, true
	}

	return template.Substitute(target.Literal, lookup)
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
