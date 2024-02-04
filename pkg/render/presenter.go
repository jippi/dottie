package render

import "github.com/jippi/dottie/pkg/ast"

type Presenter interface {
	Statement(stmt any, previous ast.Statement, settings Settings) string
	Assignment(assignment *ast.Assignment, settings Settings) string
	Comment(comment *ast.Comment, settings Settings, isAssignmentComment bool) string
	Document(doc *ast.Document, settings Settings) string
	Group(group *ast.Group, settings Settings) string
	Newline(newline *ast.Newline, settings Settings) string
	Statements(statements []ast.Statement, settings Settings) string
}
