package render

import "github.com/jippi/dottie/pkg/ast"

type Output interface {
	GroupBanner(group *ast.Group, settings Settings) *Lines
	Assignment(assignment *ast.Assignment, settings Settings) *Lines
	Comment(comment *ast.Comment, settings Settings) *Lines
	Newline(newline *ast.Newline, settings Settings) *Lines
}
