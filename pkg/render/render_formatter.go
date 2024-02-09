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
func FormatterHandler(hi *HandlerInput) HandlerSignal {
	switch statement := hi.CurrentStatement.(type) {
	case *ast.Newline:
		if !hi.Settings.showBlankLines {
			return hi.Stop()
		}

		if hi.PreviousStatement == nil {
			return hi.Return(NewLinesCollection().Newline("FormatterHandler::Newline (PreviousStatement is nil)"))
		}

		if hi.PreviousStatement.Is(&ast.Comment{}) {
			return hi.Return(NewLinesCollection().Newline("FormatterHandler::Newline (retain newlines around stand-alone comments)", hi.PreviousStatement.Type()))
		}

		// Ignore all existing newlines when doing formatting as
		// we will be injecting these ourself in other places.
		return hi.Stop()

	case *ast.Group:
		output := hi.Renderer.group(statement)
		if output.IsEmpty() {
			return hi.Stop()
		}

		buf := NewLinesCollection()

		if hi.Settings.showBlankLines && hi.PreviousStatement != nil && !hi.PreviousStatement.Is(&ast.Newline{}) {
			buf.Newline("FormatterHandler::Group:before", hi.PreviousStatement.Type())
		}

		buf.Append(output)

		return hi.Return(buf)

	case *ast.Assignment:
		output := hi.Renderer.assignment(statement)
		if output.IsEmpty() {
			return hi.Stop()
		}

		buf := NewLinesCollection()

		// If the previous Statement was also an Assignment, detect if they should
		// be allowed to cuddle (without newline between them) or not.
		//
		// Statements are only allow cuddle if both have no comments
		if hi.Settings.showBlankLines && statement.Is(hi.PreviousStatement) && (statement.HasComments() || assignmentHasComments(hi.PreviousStatement)) {
			buf.Newline("FormatterHandler::Assignment:Comments", hi.PreviousStatement.Type())
		}

		return hi.Return(buf.Append(output))
	}

	return hi.Continue()
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
