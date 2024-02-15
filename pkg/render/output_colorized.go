package render

import (
	"bytes"
	"context"

	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/tui"
)

var _ Output = (*ColorizedOutput)(nil)

type ColorizedOutput struct{}

func (ColorizedOutput) GroupBanner(ctx context.Context, group *ast.Group, settings Settings) *Lines {
	var buf bytes.Buffer

	out := tui.NewWriter(ctx, &buf).Success()

	out.Println("################################################################################")
	out.ApplyStyle(tui.Bold).Println(group.Name)
	out.Print("################################################################################")

	return NewLinesCollection().Add(buf.String())
}

func (ColorizedOutput) Assignment(ctx context.Context, assignment *ast.Assignment, settings Settings) *Lines {
	var buf bytes.Buffer

	printer := tui.NewWriter(ctx, &buf)

	if !assignment.Enabled {
		printer.Danger().Print("#")
	}

	val := assignment.Literal

	if settings.InterpolatedValues {
		val = assignment.Interpolated
	}

	printer.Primary().Print(assignment.Name)
	printer.Dark().Print("=")
	printer.Success().Print(assignment.Quote)
	printer.Warning().Print(val)
	printer.Success().Print(assignment.Quote)

	return NewLinesCollection().Add(buf.String())
}

func (ColorizedOutput) Comment(ctx context.Context, comment *ast.Comment, settings Settings) *Lines {
	var buf bytes.Buffer

	out := tui.NewWriter(ctx, &buf).Success()

	if comment.Annotation == nil {
		out.Print(comment.Value)

		return NewLinesCollection().Add(buf.String())
	}

	if comment.Annotation != nil {
		out.Print("# ")
		out.ApplyStyle(tui.Bold).Print("@", comment.Annotation.Key)
		out.Print(" ")
		out.Print(comment.Annotation.Value)
	}

	return NewLinesCollection().Add(buf.String())
}

func (ColorizedOutput) Newline(ctx context.Context, newline *ast.Newline, settings Settings) *Lines {
	if newline.Blank && !settings.ShowBlankLines() {
		return nil
	}

	return NewLinesCollection().Newline("ColorizedOutput:Newline")
}
