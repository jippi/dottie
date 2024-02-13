package upsert

import (
	"fmt"
	"slices"

	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/validation"
)

type Upserter struct {
	document              *ast.Document // The document to Upsert into
	group                 string        // The (optional) [Group] to upsert into
	placement             Placement     // The placement of the KEY in the document
	placementValue        string        // The placement value (e.g. [KEY] in [PlaceBefore] and [PlaceAfter])
	settings              Setting       // Upserter settings (bitmask)
	valuesConsideredEmpty []string      // List of values that would be considered "empty" / not-set
}

// New creates an [Upserter] with the provided settings, returning
// either the [Upserter] or an error if an [Option] validation failed
func New(document *ast.Document, options ...Option) (*Upserter, error) {
	upserter := &Upserter{
		document:  document,
		placement: AddLast,
		settings:  Validate,
	}

	if err := upserter.ApplyOptions(options...); err != nil {
		return nil, err
	}

	return upserter, nil
}

// ApplyOptions applies any additional options to the [Upserter],
// allowing you to refine and build the [Upserter] in steps.
func (u *Upserter) ApplyOptions(options ...Option) error {
	for _, option := range options {
		if err := option(u); err != nil {
			return err
		}
	}

	return nil
}

// Upsert will, depending on its options, either Update or Insert (thus, "[Up]date + In[sert]").
func (u *Upserter) Upsert(input *ast.Assignment) (*ast.Assignment, error, error) {
	assignment := u.document.Get(input.Name)
	found := assignment != nil

	// Short circuit with some quick settings checks

	switch {
	// The assignment exists, so return early
	case found && u.settings.Has(SkipIfExists):
		return nil, nil, nil

	// The assignment exists, has a literal value, and the literal value isn't what we should consider empty
	case found && u.settings.Has(SkipIfSet) && len(assignment.Literal) > 0 && !slices.Contains(u.valuesConsideredEmpty, assignment.Literal):
		return nil, nil, nil

	// The assignment exists, the literal values are the same, and they have same 'Enabled' level
	case found && u.settings.Has(SkipIfSame) && assignment.Literal == input.Literal && assignment.Enabled == input.Enabled:
		return nil, nil, nil

	// The assignment does *NOT* exists, and we require it to
	case !found && u.settings.Has(ErrorIfMissing):
		return nil, nil, fmt.Errorf("key [%s] does not exists in the document", input.Name)

	// The KEY was *NOT* found, and all other preconditions are not triggering
	case !found:
		var err error

		// Create and insert the (*ast.Assignment) into the Statement list
		assignment, err = u.createAndInsert(input)
		if err != nil {
			return nil, nil, err
		}

		// Recalculate the index order of all Statements (for interpolation)
		u.document.ReindexStatements()
	}

	// Replace comments on the assignment if the Setting is on
	if u.settings.Has(UpdateComments) {
		assignment.Comments = input.Comments
	}

	assignment.Enabled = input.Enabled
	assignment.Literal = input.Literal
	assignment.Quote = input.Quote
	assignment.Interpolated = input.Literal

	var err, warnings error

	// Interpolate the Assignment if it is enabled
	if assignment.Enabled {
		assignment.Interpolated, warnings, err = u.document.Interpolate(assignment)
		if err != nil {
			return nil, warnings, fmt.Errorf("could not interpolate variable: %w", err)
		}
	}

	// Validate
	if u.settings.Has(Validate) {
		if errors := validation.ValidateSingleAssignment(u.document, assignment.Name, nil, nil); len(errors) > 0 {
			return nil, warnings, errors[0]
		}
	}

	return assignment, warnings, nil
}

func (u *Upserter) createAndInsert(input *ast.Assignment) (*ast.Assignment, error) {
	// Ensure the group exists (may return 'nil' if no group is required)
	group := u.document.EnsureGroup(u.group)

	// Create the new newAssignment
	newAssignment := &ast.Assignment{
		Comments: input.Comments,
		Enabled:  input.Enabled,
		Group:    group,
		Literal:  input.Literal,
		Name:     input.Name,
	}

	// Find the statement slice to operate on
	statements := u.document.Statements
	if newAssignment.Group != nil {
		statements = group.Statements
	}

	var res []ast.Statement

	switch u.placement {
	// If the new assignment is desired to be first, then we prepend it to the existing
	// slice of statements
	case AddFirst:
		res = append([]ast.Statement{newAssignment}, statements...)

	// If the new assignment is desired to be last, then we append it to the existing
	// slice of statements
	case AddLast:
		res = append(statements, newAssignment)

	// If the new assignment is desired to be placed relative to another key,
	// we will figure out the ordering here
	case AddAfterKey, AddBeforeKey:
		// Run through all the statements
		for _, stmtInterface := range statements {
			// If the rangeStatement isn't an [Assignment], append it to the
			// new list of statements
			rangeStatement, ok := stmtInterface.(*ast.Assignment)
			if !ok {
				res = append(res, stmtInterface)

				continue
			}

			// If the placementValue isn't the current Assignment KEY, append it
			// to the new list of statements
			if rangeStatement.Name != u.placementValue {
				res = append(res, stmtInterface)

				continue
			}

			switch u.placement { //nolint:exhaustive
			// If placement is desired *BEFORE* another KEY, then
			//   * Append the new assignment
			//   * Append the current range statement
			case AddBeforeKey:
				res = append(res, newAssignment, rangeStatement)

			// If placement is desired *AFTER* another KEY, then
			//   * Append the current range statement
			//   * Append the new assignment
			case AddAfterKey:
				res = append(res, rangeStatement, newAssignment)
			}
		}

	default:
		// The should hopefully not happen, but just in case
		return nil, fmt.Errorf("(BUG; please report) don't know how to handle placement type: %s", u.placement)
	}

	// If the statements belonged to a Group, then update the Group with the new statement list
	if group != nil {
		group.Statements = res

		return newAssignment, nil
	}

	// Otherwise update the Document statement list
	u.document.Statements = res

	return newAssignment, nil
}
