package render

import "github.com/jippi/dottie/pkg/ast"

type Output interface {
	GroupBanner(*ast.Group, Settings) *LineBuffer
	Assignment(*ast.Assignment, Settings) *LineBuffer
	Comment(*ast.Comment, Settings) *LineBuffer
	Newline(*ast.Newline, Settings) *LineBuffer
}
