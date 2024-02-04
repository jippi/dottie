package render

import (
	"fmt"

	"github.com/jippi/dottie/pkg/ast"
)

type Renderer struct {
	Output   Outputter
	Previous ast.Statement
	Settings Settings
	handlers []Handler
}

func NewRenderer(settings Settings, additionalHandlers ...Handler) *Renderer {
	var output Outputter = Plain{}

	if settings.WithColors() {
		output = Colorized{}
	}

	handlers := append(
		[]Handler{
			FilterKeyPrefix,
			FilterActive,
			FilterGroup,
			FilterComments,
		},
		additionalHandlers...,
	)

	return &Renderer{
		Output:   output,
		Previous: nil,
		Settings: settings,
		handlers: handlers,
	}
}

func (r *Renderer) Statement(stmt any) string {
	in := &HandlerInput{
		Presenter: r,
		Previous:  r.Previous,
		Settings:  r.Settings,
		Statement: stmt,
	}

	for _, handler := range r.handlers {
		status := handler(in)

		switch status {
		// Stop processing the statement and return nothing
		case Stop:
			return ""

		// Stop processing the statement and return the value from the handler
		case Return:
			r.Previous, _ = stmt.(ast.Statement)

			return in.Value

		// Continue to next handler (or default behavior)
		case Continue:

		// Unknown signal
		default:
			panic("unknown signal: " + status.String())
		}
	}

	// Default behavior

	switch val := stmt.(type) {
	case *ast.Document:
		r.Previous = val

		return r.Document(val)

	case *ast.Group:
		r.Previous = val

		return r.Group(val)

	case *ast.Comment:
		r.Previous = val

		return r.Comment(val, false)

	case *ast.Assignment:
		r.Previous = val

		return r.Assignment(val)

	case *ast.Newline:
		r.Previous = val

		return r.Newline(val)

	case []*ast.Group:
		out := NewLineBuffer()

		for _, group := range val {
			if out.AddAndReturnPrinted(r.Statement(group)) {
				r.Previous = group
			}
		}

		return out.Get()

	case []ast.Statement:
		out := NewLineBuffer()

		for _, stmt := range val {
			if out.AddAndReturnPrinted(r.Statement(stmt)) {
				r.Previous = stmt
			}
		}

		return out.Get()

	case []*ast.Comment:
		res := NewLineBuffer()
		for _, comment := range val {
			if res.AddAndReturnPrinted(r.Statement(comment)) {
				r.Previous = comment
			}
		}

		return res.Get()

	default:
		panic(fmt.Sprintf("Unknown statement: %T", val))
	}
}

func (r *Renderer) Document(doc *ast.Document) string {
	out := NewLineBuffer()

	return out.
		Add(r.Statement(doc.Statements)).
		Add(r.Statement(doc.Groups)).
		GetWithEOF()
}

func (r *Renderer) Group(group *ast.Group) string {
	rendered := r.Statement(group.Statements)
	if len(rendered) == 0 {
		return ""
	}

	res := NewLineBuffer()

	if r.Settings.WithGroupBanners() && len(rendered) > 0 {
		res.Add(r.Output.Group(group, r.Settings))
	}

	return res.
		Add(rendered).
		Get()
}

func (r *Renderer) Assignment(a *ast.Assignment) string {
	res := NewLineBuffer()

	return res.
		Add(r.Statement(a.Comments)).
		Add(r.Output.Assignment(a, r.Settings)).
		Get()
}

func (r *Renderer) Comment(comment *ast.Comment, isAssignmentComment bool) string {
	return r.Output.Comment(comment, r.Settings, isAssignmentComment)
}

func (r *Renderer) Newline(newline *ast.Newline) string {
	return r.Output.Newline(newline, r.Settings)
}

func (r *Renderer) SetOutput(output Outputter) {
	r.Output = output
}
