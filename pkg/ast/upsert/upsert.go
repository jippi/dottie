package upsert

import (
	"errors"
	"fmt"

	"github.com/jippi/dottie/pkg/ast"
)

type Upserter struct {
	comments       []string      // Comments to add to the [Assignment]
	document       *ast.Document // The document to Upsert into
	group          string        // The (optional) [Group] to upsert into
	placement      Placement     // The placement of the KEY in the document
	placementValue string        // The placement value (e.g. [KEY] in [PlaceBefore] and [PlaceAfter])
	settings       Setting       // Upserter settings (bitmask)
}

// New creates an [Upserter] with the provided settings, returning
// either the [Upserter] or an error if an [Option] validation failed
func New(document *ast.Document, options ...Option) (*Upserter, error) {
	upserter := &Upserter{
		document:  document,
		placement: AddLast,
	}

	if err := upserter.Apply(options...); err != nil {
		return nil, err
	}

	return upserter, nil
}

func (u *Upserter) Apply(options ...Option) error {
	for _, option := range options {
		if err := option(u); err != nil {
			return err
		}
	}

	return nil
}

func (u *Upserter) Upsert(input *ast.Assignment) (*ast.Assignment, error) {
	var group *ast.Group

	existing := u.document.Get(input.Name)

	if u.settings.Has(SkipIfSet) && existing != nil && len(existing.Literal) > 0 && existing.Literal != "__CHANGE_ME__" && input.Literal != "__CHANGE_ME__" {
		return nil, nil
	}

	if u.settings.Has(SkipIfSame) && existing != nil && existing.Literal == input.Literal && existing.Active == input.Active {
		return nil, nil
	}

	found := existing != nil

	// The key does not exists!
	if !found {
		if u.settings.Has(ErrorIfMissing) {
			return nil, fmt.Errorf("Key [%s] does not exists", input.Name)
		}

		group = u.document.EnsureGroup(u.group)

		existing = &ast.Assignment{
			Name:    input.Name,
			Literal: input.Literal,
			Active:  input.Active,
			Group:   group,
		}

		existingStatements := u.document.Statements
		if existing.Group != nil {
			existingStatements = group.Statements
		}

		var res []ast.Statement

		switch u.placement {
		case AddFirst:
			res = append([]ast.Statement{existing}, existingStatements...)

		case AddLast:
			res = append(existingStatements, existing)

		case AddAfterKey, AddBeforeKey:
			for _, stmt := range existingStatements {
				assignment, ok := stmt.(*ast.Assignment)
				if !ok {
					res = append(res, stmt)

					continue
				}

				switch {
				case u.placement == AddBeforeKey && assignment.Name == u.placementValue:
					res = append(res, existing, stmt)

				case u.placement == AddAfterKey && assignment.Name == u.placementValue:
					res = append(res, stmt, existing)

				default:
					res = append(res, stmt)
				}
			}
		}

		if group != nil {
			group.Statements = res
		} else {
			u.document.Statements = res
		}
	}

	if found {
		interpolated, err := u.document.Interpolate(existing)
		if err != nil {
			return nil, errors.New("could not interpolate variable")
		}

		existing.Interpolated = interpolated
	}

	existing.Active = input.Active
	existing.Interpolated = input.Interpolated
	existing.Literal = input.Literal
	existing.Quote = input.Quote

	if comments := u.comments; len(comments) > 0 {
		existing.Comments = nil

		for _, comment := range comments {
			if len(comment) == 0 && len(comments) == 1 {
				continue
			}

			existing.Comments = append(existing.Comments, ast.NewComment(comment))
		}
	}

	return existing, nil
}
