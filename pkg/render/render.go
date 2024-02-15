package render

import (
	"context"
	"fmt"

	"github.com/jippi/dottie/pkg/ast"
)

type Renderer struct {
	Output            Output
	PreviousStatement ast.Statement
	Settings          Settings
	handlers          []Handler
}

func NewRenderer(settings Settings, additionalHandlers ...Handler) *Renderer {
	// Default handlers for filtering down the
	handlers := append(settings.Handlers(), additionalHandlers...)

	// Add Formatter handler if we're going to print pretty output!
	if settings.formatOutput {
		handlers = append(handlers, FormatterHandler)
	}

	return &Renderer{
		Output:            settings.outputter,
		PreviousStatement: nil,
		Settings:          settings,
		handlers:          handlers,
	}
}

func NewUnfilteredRenderer(settings *Settings, additionalHandlers ...Handler) *Renderer {
	return &Renderer{
		Output:            settings.outputter,
		PreviousStatement: nil,
		Settings:          *settings,
		handlers:          additionalHandlers,
	}
}

// Statement is the main loop of the Renderer.
//
// It's responsible for delegating statements to handlers, calling the right
// Output functions and track the ordering of Statements being rendered
func (r *Renderer) Statement(ctx context.Context, currentStatement any) *Lines {
	handlerInput := r.newHandlerInput(currentStatement)

	for _, handler := range r.handlers {
		status := handler(ctx, handlerInput)

		switch status {
		// Stop processing the statement and return nothing
		case Stop:
			return nil

		// Stop processing the statement and return the value from the handler
		case Return:
			if handlerInput.ReturnValue.IsEmpty() {
				return nil
			}

			// Update the "PreviousStatement" reference if
			//
			// 1) The current statement *is* a Statement (it might be a slice of statements, for example).
			// 2) The statement is *not* a group since they are rendered differently,
			//    so the statements happens "out of order" and restoring them here causes wrong values.
			if prev, ok := currentStatement.(ast.Statement); ok && !prev.Is(&ast.Group{}) {
				r.PreviousStatement = prev
			}

			return handlerInput.ReturnValue

		// Continue to next handler (or default behavior if we run out of handlers)
		case Continue:

		// Unknown signal
		default:
			panic("unknown signal: " + status.String())
		}
	}

	//
	// Default Statement behavior
	//

	switch statement := currentStatement.(type) {
	case *ast.Document:
		return r.document(ctx, statement)

	case *ast.Group:
		return r.group(ctx, statement)

	case *ast.Comment:
		return r.comment(ctx, statement)

	case *ast.Assignment:
		return r.assignment(ctx, statement)

	case *ast.Newline:
		return r.newline(ctx, statement)

	//
	// Lists of different statements will be iterated over
	//

	case []*ast.Group:
		buf := NewLinesCollection()

		for _, group := range statement {
			buf.Append(r.Statement(ctx, group))
		}

		return buf

	case []ast.Statement:
		buf := NewLinesCollection()

		for _, stmt := range statement {
			buf.Append(r.Statement(ctx, stmt))
		}

		return buf

	case []*ast.Comment:
		buf := NewLinesCollection()

		for _, comment := range statement {
			buf.Append(r.Statement(ctx, comment))
		}

		return buf

	//
	// Unrecognized Statement type
	//

	default:
		panic(fmt.Sprintf("Unknown statement: %T", statement))
	}
}

// document renders "Document" Statements.
//
// Direct Document Statements are rendered first, followed by any
// Group Statements in order they show up in the original source.
func (r *Renderer) document(ctx context.Context, document *ast.Document) *Lines {
	return NewLinesCollection().
		Append(r.Statement(ctx, document.Statements)).
		Append(r.Statement(ctx, document.Groups))
}

// group renders "Group" Statements.
func (r *Renderer) group(ctx context.Context, group *ast.Group) *Lines {
	// Capture the *current* Previous Statement in case we need to restore it (see below)
	prev := r.PreviousStatement

	// We render a Group's "Statements" before the (optional) GroupHeader (in --pretty mode).
	//
	// Because we render Group Statements "out of order" (before the Group Header),
	// we have to manually force the "Previous Statement" to be *this* Group,
	// rather than whatever *actually* was the previous statement.
	r.PreviousStatement = group

	rendered := r.Statement(ctx, group.Statements)

	if rendered.IsEmpty() {
		// If the Group Statements didn't yield any output, restore the old "PreviousStatement" before
		// any Group rendering happened, to "undo" our rendering
		r.PreviousStatement = prev

		return nil
	}

	buf := NewLinesCollection()

	// Render the optional Group banner if necessary.
	if r.Settings.ShowGroupBanners {
		buf.Append(r.Output.GroupBanner(ctx, group, r.Settings))

		if r.Settings.showBlankLines {
			buf.Newline("Group:ShowGroupBanners", r.PreviousStatement.Type(), "(type doesn't matter)")
		}
	}

	return buf.Append(rendered)
}

// assignment renders "Assignment" Statements.
func (r *Renderer) assignment(ctx context.Context, assignment *ast.Assignment) *Lines {
	// When done rendering this statement, mark it as the previous statement
	defer func() { r.PreviousStatement = assignment }()

	return NewLinesCollection().
		Append(r.Statement(ctx, assignment.Comments)).
		Append(r.Output.Assignment(ctx, assignment, r.Settings))
}

// comment renders "Comment" Statements.
func (r *Renderer) comment(ctx context.Context, comment *ast.Comment) *Lines {
	// When done rendering this statement, mark it as the previous statement
	defer func() { r.PreviousStatement = comment }()

	return r.Output.Comment(ctx, comment, r.Settings)
}

// newline renders "Newline" Statements.
func (r *Renderer) newline(ctx context.Context, newline *ast.Newline) *Lines {
	// When done rendering this statement, mark it as the previous statement
	defer func() { r.PreviousStatement = newline }()

	return r.Output.Newline(ctx, newline, r.Settings)
}

func (r *Renderer) newHandlerInput(statement any) *HandlerInput {
	return &HandlerInput{
		CurrentStatement:  statement,
		PreviousStatement: r.PreviousStatement,
		Renderer:          r,
		Settings:          r.Settings,
	}
}
