package render

import (
	"github.com/jippi/dottie/pkg/ast"
)

type Renderer struct {
	presenter *Presenter
}

func NewRenderer(presenter *Presenter) *Renderer {
	return &Renderer{
		presenter: presenter,
	}
}

func (r *Renderer) Document(doc *ast.Document, settings Settings) string {
	return r.presenter.Statement(doc, nil, settings)
}

func (r *Renderer) Group(group *ast.Group, settings Settings) string {
	return r.presenter.Statement(group, nil, settings)
}

func (r *Renderer) Assignment(assignment *ast.Assignment, settings Settings) string {
	return r.presenter.Statement(assignment, nil, settings)
}

func (r *Renderer) Comment(comment *ast.Comment, settings Settings, isAssignmentComment bool) string {
	var parent ast.Statement

	if isAssignmentComment {
		parent = &ast.Assignment{}
	}

	return r.presenter.Statement(comment, parent, settings)
}

func (r *Renderer) Newline(newline *ast.Newline, settings Settings) string {
	return r.presenter.Statement(newline, nil, settings)
}
