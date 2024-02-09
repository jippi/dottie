package render

import (
	"github.com/jippi/dottie/pkg/ast"
)

func NewFormatter() *Renderer {
	settings := Settings{
		includeDisabled:       true,
		useInterpolatedValues: false,
		showBlankLines:        true,
		showColors:            false,
		showComments:          true,
		ShowGroupBanners:      true,
		outputter:             PlainOutput{},
	}

	return NewRenderer(settings, FormatterHandler)
}

// FormatterHandler is responsible for formatting an .env file according
// to our opinionated style.
func FormatterHandler(input *HandlerInput) HandlerSignal {
	switch statement := input.CurrentStatement.(type) {
	case *ast.Newline:
		if !input.Settings.showBlankLines {
			return input.Stop()
		}

		if input.PreviousStatement == nil {
			return input.Return(NewLinesCollection().Newline("FormatterHandler::Newline (PreviousStatement is nil)"))
		}

		if input.PreviousStatement.Is(&ast.Comment{}) {
			return input.Return(NewLinesCollection().Newline("FormatterHandler::Newline (retain newlines around stand-alone comments)", input.PreviousStatement.Type()))
		}

		// Ignore all existing newlines when doing formatting as
		// we will be injecting these ourself in other places.
		return input.Stop()

	case *ast.Group:
		output := input.Renderer.group(statement)
		if output.IsEmpty() {
			return input.Stop()
		}

		buf := NewLinesCollection()

		if input.Settings.showBlankLines && input.PreviousStatement != nil && !input.PreviousStatement.Is(&ast.Newline{}) {
			buf.Newline("FormatterHandler::Group:before", input.PreviousStatement.Type())
		}

		buf.Append(output)

		return input.Return(buf)

	case *ast.Assignment:
		output := input.Renderer.assignment(statement)
		if output.IsEmpty() {
			return input.Stop()
		}

		buf := NewLinesCollection()

		// If the previous Statement was also an Assignment, detect if they should
		// be allowed to cuddle (without newline between them) or not.
		//
		// Statements are only allow cuddle if both have no comments
		if input.Settings.showBlankLines && statement.Is(input.PreviousStatement) && (statement.HasComments() || assignmentHasComments(input.PreviousStatement)) {
			buf.Newline("FormatterHandler::Assignment:Comments", input.PreviousStatement.Type())
		}

		return input.Return(buf.Append(output))
	}

	return input.Continue()
}

// assignmentHasComments checks if the Statement is an Assignment
// and if it has any comments attached to it
func assignmentHasComments(statement ast.Statement) bool {
	assignment, ok := statement.(*ast.Assignment)
	if !ok {
		return false
	}

	return assignment.HasComments()
}
