package render

import (
	"github.com/jippi/dottie/pkg/ast"
)

var _ Output = (*CompletionOutputKeys)(nil)

type CompletionOutputKeys struct{}

func (CompletionOutputKeys) GroupBanner(group *ast.Group, settings Settings) *Lines {
	return nil
}

func (CompletionOutputKeys) Assignment(a *ast.Assignment, settings Settings) *Lines {
	return NewLinesCollection().Add(a.Name)
}

func (r CompletionOutputKeys) Comment(comment *ast.Comment, settings Settings) *Lines {
	return nil
}

func (r CompletionOutputKeys) Newline(newline *ast.Newline, settings Settings) *Lines {
	return nil
}
