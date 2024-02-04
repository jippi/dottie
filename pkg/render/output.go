package render

import "github.com/jippi/dottie/pkg/ast"

type Output interface {
	GroupBanner(*ast.Group, Settings) string
	Assignment(*ast.Assignment, Settings) string
	Comment(*ast.Comment, Settings) string
	Newline(*ast.Newline, Settings) string
}
