package render

import (
	"context"

	"github.com/jippi/dottie/pkg/ast"
)

type Output interface {
	GroupBanner(ctx context.Context, group *ast.Group, settings Settings) *Lines
	Assignment(ctx context.Context, assignment *ast.Assignment, settings Settings) *Lines
	Comment(ctx context.Context, comment *ast.Comment, settings Settings) *Lines
	Newline(ctx context.Context, newline *ast.Newline, settings Settings) *Lines
}
