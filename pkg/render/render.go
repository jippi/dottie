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

func (r *Renderer) Statement(currentStatement any) string {
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
			return ""

		// Stop processing the statement and return the value from the handler
		case Return:
			r.PreviousStatement, _ = currentStatement.(ast.Statement)

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
		r.PreviousStatement = statement

		return r.Document(statement)

	case *ast.Group:
		r.PreviousStatement = statement

		return r.Group(statement)

	case *ast.Comment:
		r.PreviousStatement = statement

		return r.Comment(statement)

	case *ast.Assignment:
		r.PreviousStatement = statement

		return r.Assignment(statement)

	case *ast.Newline:
		r.PreviousStatement = statement

		return r.Newline(statement)

	//
	// Lists of different statements will be iterated over
	//

	case []*ast.Group:
		buf := NewLineBuffer()

		for _, group := range statement {
			if buf.AddAndReturnPrinted(r.Statement(group)) {
				r.PreviousStatement = group
			}
		}

		return buf.Get()

	case []ast.Statement:
		buf := NewLineBuffer()

		for _, stmt := range statement {
			if buf.AddAndReturnPrinted(r.Statement(stmt)) {
				r.PreviousStatement = stmt
			}
		}

		return buf.Get()

	case []*ast.Comment:
		buf := NewLineBuffer()
		for _, comment := range statement {
			if buf.AddAndReturnPrinted(r.Statement(comment)) {
				r.PreviousStatement = comment
			}
		}

		return buf.Get()

	//
	// Unrecognized Statement type
	//

	default:
		panic(fmt.Sprintf("Unknown statement: %T", statement))
	}
}

func (r *Renderer) Document(document *ast.Document) string {
	return NewLineBuffer().
		Add(r.Statement(document.Statements)).
		Add(r.Statement(document.Groups)).
		GetWithEOF()
}

func (r *Renderer) Group(group *ast.Group) string {
	rendered := r.Statement(group.Statements)
	if len(rendered) == 0 {
		return ""
	}

	buf := NewLineBuffer()

	if r.Settings.ShowGroupBanners && len(rendered) > 0 {
		buf.Add(r.Output.GroupBanner(group, r.Settings))
	}

	return buf.Add(rendered).Get()
}

func (r *Renderer) Assignment(assignment *ast.Assignment) string {
	return NewLineBuffer().
		Add(r.Statement(assignment.Comments)).
		Add(r.Output.Assignment(assignment, r.Settings)).
		Get()
}

func (r *Renderer) Comment(comment *ast.Comment) string {
	return r.Output.Comment(comment, r.Settings)
}

func (r *Renderer) Newline(newline *ast.Newline) string {
	return r.Output.Newline(newline, r.Settings)
}
