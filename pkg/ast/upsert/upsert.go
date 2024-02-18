package upsert

import (
	"context"
	"fmt"
	"slices"

	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/parser"
	"github.com/jippi/dottie/pkg/render"
	"github.com/jippi/dottie/pkg/scanner"
	"go.uber.org/multierr"
)

type Upserter struct {
	document              *ast.Document // The document to Upsert into
	group                 string        // The (optional) [Group] to upsert into
	placement             Placement     // The placement of the KEY in the document
	placementValue        string        // The placement value (e.g. [KEY] in [PlaceBefore] and [PlaceAfter])
	settings              Setting       // Upserter settings (bitmask)
	valuesConsideredEmpty []string      // List of values that would be considered "empty" / not-set
	ignoreValidationRules []string      // Validation rules that should be ignored
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
func (u *Upserter) Upsert(ctx context.Context, input *ast.Assignment) (*ast.Assignment, error, error) {
	assignment := u.document.Get(input.Name)
	exists := assignment != nil

	// Short circuit with some quick settings checks

	switch {
	// The assignment exists, so return early
	case exists && u.settings.Has(SkipIfExists):
		return nil, nil, nil

	// The assignment does *NOT* exists, and we require it to
	case !exists && u.settings.Has(ErrorIfMissing):
		return nil, nil, fmt.Errorf("key [%s] does not exists in the document", input.Name)

	// The assignment exists, has a literal value, and the literal value isn't what we should consider empty
	case exists && u.settings.Has(SkipIfSet) && len(assignment.Literal) > 0 && !slices.Contains(u.valuesConsideredEmpty, assignment.Literal):
		return nil, nil, nil

	// The assignment exists, the literal values are the same, and they have same 'Enabled' level
	case exists && u.settings.Has(SkipIfSame) && assignment.Literal == input.Literal && assignment.Enabled == input.Enabled:
		return nil, nil, nil

	// The KEY was *NOT* found, and all other preconditions are not triggering
	case !exists:
		var err error

		// Create and insert the (*ast.Assignment) into the Statement list
		assignment, err = u.createAndInsert(ctx, input)
		if err != nil {
			return nil, nil, err
		}

		// Make sure to reindex the document
		u.document.ReindexStatements()
	}

	// Replace comments on the assignment if the Setting is on
	if u.settings.Has(UpdateComments) {
		assignment.Comments = input.Comments
	}

	assignment.Enabled = input.Enabled
	assignment.Literal = input.Literal
	assignment.Interpolated = input.Literal
	assignment.Quote = input.Quote

	var (
		tempDoc       *ast.Document
		err, warnings error
	)

	// Render and parse back the Statement to ensure annotations and such are properly handled
	thing := u.document.AllAssignments()[:assignment.Position.Index+1]

	tempDoc, err = parser.New(scanner.New(render.NewFormatter().Statement(ctx, thing).String()), "memory://tmp").Parse()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse assignment: %w", err)
	}

	assignment = tempDoc.Get(assignment.Name)
	assignment.Initialize()

	if _, ok := assignment.Dependencies[assignment.Name]; ok {
		return nil, nil, fmt.Errorf("Key [%s] may not reference itself!", assignment.Name)
	}

	// Replace the Assignment in the document
	//
	// This is necessary since its a different pointer address after we rendered+parsed earlier
	u.document.Replace(assignment)

	// Reinitialize the document so all indices and such are correct
	u.document.Initialize()

	// Interpolate the Assignment if it is enabled
	if assignment.Enabled {
		warnings, err = u.document.InterpolateStatement(assignment)
		if err != nil {
			return nil, warnings, err
		}
	}

	// Validate
	if u.settings.Has(Validate) {
		if validationErrors, warns, errs := u.document.ValidateSingleAssignment(assignment, nil, u.ignoreValidationRules); len(validationErrors) > 0 {
			warnings = multierr.Append(warnings, warns)
			errs = multierr.Append(errs, validationErrors)

			return assignment, warnings, errs
		}
	}

	return assignment, warnings, nil
}

func (u *Upserter) createAndInsert(ctx context.Context, input *ast.Assignment) (*ast.Assignment, error) {
	// Create the new newAssignment
	newAssignment := &ast.Assignment{
		Comments: input.Comments,
		Enabled:  input.Enabled,
		Literal:  input.Literal,
		Name:     input.Name,
	}

	doc, err := parser.New(scanner.New(render.NewFormatter().Statement(ctx, newAssignment).String()), "-").Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse assignment: %w", err)
	}

	// Ensure the group exists (may return 'nil' if no group is required)
	group := u.document.EnsureGroup(u.group)

	newAssignment = doc.Get(newAssignment.Name)
	newAssignment.Group = group

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
