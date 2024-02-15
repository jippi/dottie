package render

import (
	"context"

	"github.com/jippi/dottie/pkg/ast"
)

var _ Output = (*CompletionOutputKeys)(nil)

type CompletionOutputKeys struct{}

func (CompletionOutputKeys) GroupBanner(ctx context.Context, group *ast.Group, settings Settings) *Lines {
	return nil
}

func (CompletionOutputKeys) Assignment(ctx context.Context, assignment *ast.Assignment, settings Settings) *Lines {
	return NewLinesCollection().Add(assignment.Name + "\t" + assignment.DocumentationSummary())
}

func (r CompletionOutputKeys) Comment(ctx context.Context, comment *ast.Comment, settings Settings) *Lines {
	return nil
}

func (r CompletionOutputKeys) Newline(ctx context.Context, newline *ast.Newline, settings Settings) *Lines {
	return nil
}
