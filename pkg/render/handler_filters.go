package render

import (
	"strings"

	"github.com/jippi/dottie/pkg/ast"
)

func FilterKeyPrefix(in *HandlerInput) HandlerSignal {
	switch val := in.Statement.(type) {
	case *ast.Assignment:
		if len(in.Settings.FilterGroup) > 0 && !strings.HasPrefix(val.Name, in.Settings.FilterKeyPrefix) {
			return in.Stop()
		}
	}

	return in.Continue()
}

func FilterComments(in *HandlerInput) HandlerSignal {
	switch in.Statement.(type) {
	case *ast.Comment:
		if !in.Settings.WithComments() {
			return in.Stop()
		}
	}

	return in.Continue()
}

func FilterActive(in *HandlerInput) HandlerSignal {
	switch val := in.Statement.(type) {
	case *ast.Assignment:
		if !val.Active && !in.Settings.IncludeCommented {
			return in.Stop()
		}
	}

	return in.Continue()
}

func FilterGroup(in *HandlerInput) HandlerSignal {
	switch val := in.Statement.(type) {
	case *ast.Assignment:
		if !val.BelongsToGroup(in.Settings.FilterGroup) {
			return in.Stop()
		}

	case *ast.Group:
		if !val.BelongsToGroup(in.Settings.FilterGroup) {
			return in.Stop()
		}

	case *ast.Comment:
		if !val.BelongsToGroup(in.Settings.FilterGroup) {
			return in.Stop()
		}
	}

	return in.Continue()
}
