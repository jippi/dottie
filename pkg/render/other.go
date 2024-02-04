package render

import "github.com/jippi/dottie/pkg/ast"

func RenderFull(doc *ast.Document) string {
	return (&PlainPresenter{}).Document(doc, Settings{
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
