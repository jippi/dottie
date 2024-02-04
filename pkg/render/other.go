package render

import "github.com/jippi/dottie/pkg/ast"

func NewFormatter() *Renderer {
	settings := Settings{
		IncludeCommented: true,
		Interpolate:      false,
		ShowBlankLines:   true,
		ShowColors:       false,
		ShowComments:     true,
		ShowGroups:       true,
	}

	return NewRenderer(settings, Format)
}

func assignmentHasComments(stmt ast.Statement) bool {
	x, ok := stmt.(*ast.Assignment)
	if !ok {
		return false
	}

	return x.HasComments()
}
