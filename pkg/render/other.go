package render

import "github.com/jippi/dottie/pkg/ast"

func RenderFormatted(doc *ast.Document) string {
	return NewRenderer(&FormattedPresenter{}).
		Document(
			doc,
			Settings{
				IncludeCommented: true,
				Interpolate:      false,
				ShowBlankLines:   true,
				ShowColors:       false,
				ShowComments:     true,
				ShowGroups:       true,
			},
		)
}

func RenderDirect(doc *ast.Document) string {
	return NewRenderer(&DirectPresenter{}).
		Document(
			doc,
			Settings{
				IncludeCommented: true,
				Interpolate:      false,
				ShowBlankLines:   true,
				ShowColors:       false,
				ShowComments:     true,
				ShowGroups:       true,
			},
		)
}

func assignmentHasComments(stmt ast.Statement) bool {
	x, ok := stmt.(*ast.Assignment)
	if !ok {
		return false
	}

	return x.HasComments()
}
