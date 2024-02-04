package render

import (
	"fmt"

	"github.com/jippi/dottie/pkg/ast"
)

type DirectPresenter struct {
	output Outputter
}

var _ Presenter = (*DirectPresenter)(nil)

func NewDirectPresenter(settings Settings) *DirectPresenter {
	var output Outputter = Plain{}

	if settings.WithColors() {
		output = Colorized{}
	}

	return &DirectPresenter{
		output: output,
	}
}

func (r *DirectPresenter) Statement(stmt any, previous ast.Statement, settings Settings) string {
	switch val := stmt.(type) {
	case *ast.Document:
		return r.Document(val, settings)

	case *ast.Group:
		return r.Group(val, settings)

	case *ast.Comment:
		return r.Comment(val, settings, false)

	case *ast.Assignment:
		return r.Assignment(val, settings)

	case *ast.Newline:
		return r.Newline(val, settings)

	case []*ast.Group:
		var (
			out  = &LineBuffer{}
			prev ast.Statement
		)

		for _, x := range val {
			if out.AddPrinted(r.Statement(x, prev, settings)) {
				prev = x
			}
		}

		return out.Get()

	case []ast.Statement:
		var (
			out  = &LineBuffer{}
			prev ast.Statement
		)

		for _, stmt := range val {
			if out.AddPrinted(r.Statement(stmt, prev, settings)) {
				prev = stmt
			}
		}

		return out.Get()

	case []*ast.Comment:
		if !settings.WithComments() {
			return ""
		}

		res := LineBuffer{}
		for _, c := range val {
			res.Add(r.Comment(c, settings, true))
		}

		return res.Get()

	default:
		panic(fmt.Sprintf("Unknown statement: %T", val))
	}
}

func (r *DirectPresenter) Document(doc *ast.Document, settings Settings) string {
	out := &LineBuffer{}

	return out.
		Add(r.Statement(doc.Statements, doc, settings)).
		Add(r.Statement(doc.Groups, doc, settings)).
		Get()
}

func (r *DirectPresenter) Group(group *ast.Group, settings Settings) string {
	if !group.BelongsToGroup(settings.FilterGroup) {
		return ""
	}

	rendered := r.Statement(group.Statements, group, settings)
	if len(rendered) == 0 {
		return ""
	}

	res := &LineBuffer{}

	if settings.WithGroups() && len(rendered) > 0 {
		res.Add(r.output.Group(group, settings))
	}

	return res.
		Add(rendered).
		Get()
}

func (r *DirectPresenter) Assignment(a *ast.Assignment, settings Settings) string {
	if !settings.Match(a) || !a.BelongsToGroup(settings.FilterGroup) {
		return ""
	}

	res := &LineBuffer{}

	return res.
		Add(r.Statement(a.Comments, a, settings)).
		Add(r.output.Assignment(a, settings)).
		Get()
}

func (r *DirectPresenter) Comment(comment *ast.Comment, settings Settings, isAssignmentComment bool) string {
	if !settings.WithComments() {
		return ""
	}

	return r.output.Comment(comment, settings, isAssignmentComment)
}

func (r *DirectPresenter) Newline(newline *ast.Newline, settings Settings) string {
	return r.output.Newline(newline, settings)
}

func (r *DirectPresenter) SetOutput(output Outputter) {
	r.output = output
}
