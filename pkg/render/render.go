package render

import (
	"github.com/jippi/dottie/pkg/ast"
)

type Renderer interface {
	Document(doc *ast.Document, settings Settings) string
	Group(group *ast.Group, settings Settings) string
	Assignment(assignment *ast.Assignment, settings Settings) string
	Comment(comment *ast.Comment, settings Settings) string
	Newline(newline *ast.Newline, settings Settings) string
}

func RenderFull(doc *ast.Document) string {
	return (&PlainRenderer{}).Document(doc, Settings{
		IncludeCommented: true,
		Interpolate:      false,
		ShowBlankLines:   true,
		ShowColors:       false,
		ShowComments:     true,
		ShowGroups:       true,
	})
}

func assignmentHasComments(stmt ast.Statement) bool {
	x, ok := stmt.(*ast.Assignment)
	if !ok {
		return false
	}

	return x.HasComments()
}
