package render

import "github.com/jippi/dottie/pkg/ast"

func NewFormatter(doc *ast.Document) string {
	settings := Settings{
		IncludeCommented: true,
		Interpolate:      false,
		ShowBlankLines:   true,
		ShowColors:       false,
		ShowComments:     true,
		ShowGroups:       true,
	}

	return NewRenderer(settings, Format).
		Document(doc, settings)
}

func NewDirect(doc *ast.Document) string {
	settings := Settings{
		IncludeCommented: true,
		Interpolate:      false,
		ShowBlankLines:   true,
		ShowColors:       false,
		ShowComments:     true,
		ShowGroups:       true,
	}

	return NewRenderer(settings).
		Document(doc, settings)
}

func assignmentHasComments(stmt ast.Statement) bool {
	x, ok := stmt.(*ast.Assignment)
	if !ok {
		return false
	}

	return x.HasComments()
}
