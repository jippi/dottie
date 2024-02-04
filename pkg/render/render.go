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

func NewRenderer(settings Settings, handlers ...Handler) *Renderer {
	var output Outputter = Plain{}

	if settings.WithColors() {
		output = Colorized{}
	}

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
		out := &LineBuffer{}

		for _, group := range val {
			if out.AddPrinted(r.Statement(group)) {
				r.Previous = group
			}
		}

		return out.Get()

	case []ast.Statement:
		out := &LineBuffer{}

		for _, stmt := range val {
			if out.AddPrinted(r.Statement(stmt)) {
				r.Previous = stmt
			}
		}

		return out.Get()

	case []*ast.Comment:
		if !r.Settings.WithComments() {
			return ""
		}

		res := LineBuffer{}
		for _, comment := range val {
			if res.AddPrinted(r.Comment(comment, true)) {
				r.Previous = comment
			}
		}

		return res.Get()

	default:
		panic(fmt.Sprintf("Unknown statement: %T", val))
	}
}

func (r *Renderer) Document(doc *ast.Document) string {
	out := &LineBuffer{}

	return out.
		Add(r.Statement(doc.Statements)).
		Add(r.Statement(doc.Groups)).
		GetWithEOF()
}

func (r *Renderer) Group(group *ast.Group) string {
	if !group.BelongsToGroup(r.Settings.FilterGroup) {
		return ""
	}

	rendered := r.Statement(group.Statements)
	if len(rendered) == 0 {
		return ""
	}

	res := &LineBuffer{}

	if r.Settings.WithGroups() && len(rendered) > 0 {
		res.Add(r.Output.Group(group, r.Settings))
	}

	return res.
		Add(rendered).
		Get()
}

func (r *Renderer) Assignment(a *ast.Assignment) string {
	if !r.Settings.Match(a) || !a.BelongsToGroup(r.Settings.FilterGroup) {
		return ""
	}

	res := &LineBuffer{}

	return res.
		Add(r.Statement(a.Comments)).
		Add(r.Output.Assignment(a, r.Settings)).
		Get()
}

func (r *Renderer) Comment(comment *ast.Comment, isAssignmentComment bool) string {
	if !r.Settings.WithComments() {
		return ""
	}

	return r.Output.Comment(comment, r.Settings, isAssignmentComment)
}

func (r *Renderer) Newline(newline *ast.Newline) string {
	return r.Output.Newline(newline, r.Settings)
}

func (r *Renderer) SetOutput(output Outputter) {
	r.Output = output
}
