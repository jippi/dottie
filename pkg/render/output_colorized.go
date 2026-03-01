package render

import (
	"context"
	"strings"

	"github.com/jippi/dottie/pkg/ast"
	"github.com/jippi/dottie/pkg/tui"
)

var _ Output = (*ColorizedOutput)(nil)

type ColorizedOutput struct{}

func (ColorizedOutput) GroupBanner(ctx context.Context, group *ast.Group, settings Settings) *Lines {
	writer := tui.NewWriter(ctx, nil)
	success := writer.Success()

	return NewLinesCollection().
		Add(success.Sprint("################################################################################")).
		Add(success.ApplyStyle(tui.Bold).Sprint(group.Name)).
		Add(success.Sprint("################################################################################"))
}

func (ColorizedOutput) Assignment(ctx context.Context, assignment *ast.Assignment, settings Settings) *Lines {
	printer := tui.NewWriter(ctx, nil)

	var out strings.Builder

	if !assignment.Enabled {
		out.WriteString(printer.Danger().Sprint("#"))
	}

	if settings.export {
		out.WriteString(printer.Dark().Sprint("export "))
	}

	val := assignment.Literal

	if settings.InterpolatedValues {
		val = assignment.Interpolated
	}

	out.WriteString(printer.Primary().Sprint(assignment.Name))
	out.WriteString(printer.Dark().Sprint("="))
	out.WriteString(printer.Success().Sprint(assignment.Quote))
	out.WriteString(printer.Warning().Sprint(val))
	out.WriteString(printer.Success().Sprint(assignment.Quote))

	return NewLinesCollection().Add(out.String())
}

func (ColorizedOutput) Comment(ctx context.Context, comment *ast.Comment, settings Settings) *Lines {
	writer := tui.NewWriter(ctx, nil)
	out := writer.Success()

	if comment.Annotation == nil {
		return NewLinesCollection().Add(out.Sprint(comment.Value))
	}

	var builder strings.Builder

	builder.WriteString(out.Sprint("# "))
	builder.WriteString(out.ApplyStyle(tui.Bold).Sprint("@", comment.Annotation.Key))
	builder.WriteString(out.Sprint(" "))
	builder.WriteString(out.Sprint(comment.Annotation.Value))

	return NewLinesCollection().Add(builder.String())
}

func (ColorizedOutput) Newline(ctx context.Context, newline *ast.Newline, settings Settings) *Lines {
	if newline.Blank && !settings.ShowBlankLines() {
		return nil
	}

	return NewLinesCollection().Newline("ColorizedOutput:Newline")
}
