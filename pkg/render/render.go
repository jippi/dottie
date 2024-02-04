package render

import (
	"fmt"

	"github.com/jippi/dottie/pkg/ast"
)

type Renderer struct {
	Output   Outputter
	Previous ast.Statement
	handlers []Handler
}

func NewRenderer(settings Settings, handlers ...Handler) *Renderer {
	var output Outputter = Plain{}

	if settings.WithColors() {
		output = Colorized{}
	}

	return &Renderer{
		Output:   output,
		handlers: handlers,
	}
}

func (r *Renderer) Statement(stmt any, settings Settings) string {
	in := &HandlerInput{
		Presenter: r,
		Previous:  r.Previous,
		Settings:  settings,
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

		return r.Document(val, settings)

	case *ast.Group:
		r.Previous = val

		return r.Group(val, settings)

	case *ast.Comment:
		r.Previous = val

		return r.Comment(val, settings, false)

	case *ast.Assignment:
		r.Previous = val

		return r.Assignment(val, settings)

	case *ast.Newline:
		r.Previous = val

		return r.Newline(val, settings)

	case []*ast.Group:
		out := &LineBuffer{}

		for _, group := range val {
			if out.AddPrinted(r.Statement(group, settings)) {
				r.Previous = group
			}
		}

		return out.Get()

	case []ast.Statement:
		out := &LineBuffer{}

		for _, stmt := range val {
			if out.AddPrinted(r.Statement(stmt, settings)) {
				r.Previous = stmt
			}
		}

		return out.Get()

	case []*ast.Comment:
		if !settings.WithComments() {
			return ""
		}

		res := LineBuffer{}
		for _, comment := range val {
			if res.AddPrinted(r.Comment(comment, settings, true)) {
				r.Previous = comment
			}
		}

		return res.Get()

	default:
		panic(fmt.Sprintf("Unknown statement: %T", val))
	}
}

func (r *Renderer) Document(doc *ast.Document, settings Settings) string {
	out := &LineBuffer{}

	return out.
		Add(r.Statement(doc.Statements, settings)).
		Add(r.Statement(doc.Groups, settings)).
		Get()
}

func (r *Renderer) Group(group *ast.Group, settings Settings) string {
	if !group.BelongsToGroup(settings.FilterGroup) {
		return ""
	}

	rendered := r.Statement(group.Statements, settings)
	if len(rendered) == 0 {
		return ""
	}

	res := &LineBuffer{}

	if settings.WithGroups() && len(rendered) > 0 {
		res.Add(r.Output.Group(group, settings))
	}

	return res.
		Add(rendered).
		Get()
}

func (r *Renderer) Assignment(a *ast.Assignment, settings Settings) string {
	if !settings.Match(a) || !a.BelongsToGroup(settings.FilterGroup) {
		return ""
	}

	res := &LineBuffer{}

	return res.
		Add(r.Statement(a.Comments, settings)).
		Add(r.Output.Assignment(a, settings)).
		Get()
}

func (r *Renderer) Comment(comment *ast.Comment, settings Settings, isAssignmentComment bool) string {
	if !settings.WithComments() {
		return ""
	}

	return r.Output.Comment(comment, settings, isAssignmentComment)
}

func (r *Renderer) Newline(newline *ast.Newline, settings Settings) string {
	return r.Output.Newline(newline, settings)
}

func (r *Renderer) SetOutput(output Outputter) {
	r.Output = output
}
