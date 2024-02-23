package render

import (
	"bytes"
	"context"
	"fmt"

	"github.com/jippi/dottie/pkg/ast"
)

var _ Output = (*PlainOutput)(nil)

type PlainOutput struct{}

func (PlainOutput) GroupBanner(ctx context.Context, group *ast.Group, settings Settings) *Lines {
	out := NewLinesCollection()

	out.Add("################################################################################")
	out.Add(group.Name)
	out.Add("################################################################################")

	return out
}

func (PlainOutput) Assignment(ctx context.Context, assignment *ast.Assignment, settings Settings) *Lines {
	var buf bytes.Buffer

	if !assignment.Enabled {
		buf.WriteString("#")
	}

	val := assignment.Literal

	if settings.InterpolatedValues {
		val = assignment.Interpolated
	}

	fmt.Println("YOOO >", fmt.Sprintf("%s=%s%s%s", assignment.Name, assignment.Quote, val, assignment.Quote), "<")
	buf.WriteString(fmt.Sprintf("%s=%s%s%s", assignment.Name, assignment.Quote, val, assignment.Quote))

	return NewLinesCollection().Add(buf.String())
}

func (r PlainOutput) Comment(ctx context.Context, comment *ast.Comment, settings Settings) *Lines {
	return NewLinesCollection().Add(comment.Value)
}

func (r PlainOutput) Newline(ctx context.Context, newline *ast.Newline, settings Settings) *Lines {
	if newline.Blank && !settings.ShowBlankLines() {
		return nil
	}

	return NewLinesCollection().Newline("PlainOutput:Newline")
}
