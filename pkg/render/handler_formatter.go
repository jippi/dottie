package render

import (
	"github.com/jippi/dottie/pkg/ast"
)

func Format(in *HandlerInput) Signal {
	switch val := in.Statement.(type) {
	// Ignore all existing newlines when doing formatting
	// we will be injecting these ourself in other places
	case *ast.Newline:
		return in.Stop()

	case *ast.Group:
		output := in.Presenter.Group(val, in.Settings)
		if len(output) == 0 {
			return in.Stop()
		}

		res := &LineBuffer{}

		// If the previous line is a newline, don't add another one.
		// This could happen if a group is the *first* thing in the document
		if !(&ast.Newline{}).Is(in.Previous) && in.Previous != nil {
			res.Newline()
		}

		res.Add(output)
		res.Newline()

		return in.Return(res.Get())

	case *ast.Assignment:
		output := in.Presenter.Assignment(val, in.Settings)
		if len(output) == 0 {
			return in.Stop()
		}

		buff := LineBuffer{}

		// If the assignment belongs to a group, but there are no previous
		// then we're the first, so add a newline padding
		if val.Group != nil && in.Previous == nil {
			buff.Newline()
		}

		// Looks like current and previous Statement is both "Assignment"
		// which mean they might be too close in the document, so we will
		// attempt to inject some new-lines to give them some space
		if in.Settings.ShowPretty && val.Is(in.Previous) {
			// only allow cuddling of assignments if they both have no comments
			if val.HasComments() || assignmentHasComments(in.Previous) {
				buff.Newline()
			}
		}

		return in.Return(buff.Add(output).Get())
	}

	return in.Continue()
}
