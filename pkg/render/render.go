package render

import (
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
	var output Output = PlainOutput{}

	if settings.ShowColors {
		output = ColorizedOutput{}
	}

	// Default handlers for filtering down the
	handlers := append(
		[]Handler{
			FilterDisabledStatements,
			FilterByKeyPrefix,
			FilterByGroupName,
			FilterComments,
		},
		additionalHandlers...,
	)

	return &Renderer{
		Output:            output,
		PreviousStatement: nil,
		Settings:          settings,
		handlers:          handlers,
	}
}

func (r *Renderer) Statement(currentStatement any) *LineBuffer {
	hi := &HandlerInput{
		Presenter:         r,
		PreviousStatement: r.PreviousStatement,
		Settings:          r.Settings,
		CurrentStatement:  currentStatement,
	}

	for _, handler := range r.handlers {
		status := handler(hi)

		switch status {
		// Stop processing the statement and return nothing
		case Stop:
			return nil

		// Stop processing the statement and return the value from the handler
		case Return:
			if prev, ok := currentStatement.(ast.Statement); ok && !hi.ReturnValue.Empty() && !prev.Is(&ast.Group{}) {
				r.PreviousStatement = prev
			}

			return hi.ReturnValue

		// Continue to next handler (or default behavior)
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
		return r.Document(statement)

	case *ast.Group:
		return r.Group(statement)

	case *ast.Comment:
		return r.Comment(statement)

	case *ast.Assignment:
		return r.Assignment(statement)

	case *ast.Newline:
		return r.Newline(statement)

	//
	// Lists of different statements will be iterated over
	//

	case []*ast.Group:
		buf := NewLineBuffer()

		for _, group := range statement {
			buf.Add(r.Statement(group))
		}

		return buf

	case []ast.Statement:
		buf := NewLineBuffer()

		for _, stmt := range statement {
			buf.Add(r.Statement(stmt))
		}

		return buf

	case []*ast.Comment:
		buf := NewLineBuffer()

		for _, comment := range statement {
			buf.Add(r.Statement(comment))
		}

		return buf

	//
	// Unrecognized Statement type
	//

	default:
		panic(fmt.Sprintf("Unknown statement: %T", statement))
	}
}

func (r *Renderer) Document(document *ast.Document) *LineBuffer {
	return NewLineBuffer().
		Add(r.Statement(document.Statements)).
		Add(r.Statement(document.Groups))
}

func (r *Renderer) Group(group *ast.Group) *LineBuffer {
	prev := r.PreviousStatement

	// Render groups inner statements with the group being "previous"
	// This is necessary because we render the Group statements *before* the (optional)
	// GroupHeader, so for things to detect and align itself correctly, we need to fake the behavior of rendering order
	r.PreviousStatement = group

	rendered := r.Statement(group.Statements)
	if rendered.Empty() {
		r.PreviousStatement = prev

		return nil
	}

	buf := NewLineBuffer()

	if r.Settings.ShowGroupBanners {
		buf.
			Add(r.Output.GroupBanner(group, r.Settings)).
			AddNewline("Group:ShowGroupBanners", r.PreviousStatement.Type(), "(type doesn't matter)")
	}

	return buf.Add(rendered)
}

func (r *Renderer) Assignment(assignment *ast.Assignment) *LineBuffer {
	defer func() { r.PreviousStatement = assignment }()

	return NewLineBuffer().
		Add(r.Statement(assignment.Comments)).
		Add(r.Output.Assignment(assignment, r.Settings))
}

func (r *Renderer) Comment(comment *ast.Comment) *LineBuffer {
	defer func() { r.PreviousStatement = comment }()

	return r.Output.Comment(comment, r.Settings)
}

func (r *Renderer) Newline(newline *ast.Newline) *LineBuffer {
	defer func() { r.PreviousStatement = newline }()

	return r.Output.Newline(newline, r.Settings)
}
