package render

import (
	"strings"

	"github.com/jippi/dottie/pkg/ast"
)

// FilterKeyPrefix will filter out Statements that do not have the
// configured (optional) key prefix
func FilterKeyPrefix(in *HandlerInput) HandlerSignal {
	// Short circuit the filter if there is no KeyPrefix to filter on
	if len(in.Settings.FilterKeyPrefix) == 0 {
		return in.Continue()
	}

	switch val := in.Statement.(type) {
	case *ast.Assignment:
		if !strings.HasPrefix(val.Name, in.Settings.FilterKeyPrefix) {
			return in.Stop()
		}
	}

	return in.Continue()
}

// FilterComments will filter out Comment statements if they aren't to be included
func FilterComments(in *HandlerInput) HandlerSignal {
	// Short circuit the filter if we allow comments
	if in.Settings.WithComments() {
		return in.Continue()
	}

	switch in.Statement.(type) {
	case *ast.Comment:
		if !in.Settings.WithComments() {
			return in.Stop()
		}
	}

	return in.Continue()
}

// FilterDisabledStatements will filter out Assignment Statements that are
// disabled
func FilterDisabledStatements(in *HandlerInput) HandlerSignal {
	// Short circuit the filter if we allow disabled statements
	if in.Settings.IncludeDisabled {
		return in.Continue()
	}

	switch val := in.Statement.(type) {
	case *ast.Assignment:
		if !val.Active && !in.Settings.IncludeDisabled {
			return in.Stop()
		}
	}

	return in.Continue()
}

// FilterGroupName will filter out Statements that do not
// belong to the required Group name
func FilterGroupName(in *HandlerInput) HandlerSignal {
	// Short circuit the filter if there is no Group name to filter on
	if len(in.Settings.FilterGroup) == 0 {
		return in.Continue()
	}

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
