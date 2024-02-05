package render

import "github.com/jippi/dottie/pkg/ast"

type Output interface {
	GroupBanner(*ast.Group, Settings) *Lines
	Assignment(*ast.Assignment, Settings) *Lines
	Comment(*ast.Comment, Settings) *Lines
	Newline(*ast.Newline, Settings) *Lines
}
