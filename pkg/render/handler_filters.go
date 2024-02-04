package render

import (
	"strings"

	"github.com/jippi/dottie/pkg/ast"
)

// FilterByKeyPrefix will filter out Statements that do not have the
// configured (optional) key prefix
func FilterByKeyPrefix(hi *HandlerInput) HandlerSignal {
	// Short circuit the filter if there is no KeyPrefix to filter on
	if len(hi.Settings.FilterKeyPrefix) == 0 {
		return hi.Continue()
	}

	switch statement := hi.CurrentStatement.(type) {
	case *ast.Assignment:
		if !strings.HasPrefix(statement.Name, hi.Settings.FilterKeyPrefix) {
			return hi.Stop()
		}
	}

	return hi.Continue()
}

// FilterComments will filter out Comment statements if they aren't to be included
func FilterComments(hi *HandlerInput) HandlerSignal {
	// Short circuit the filter if we allow comments
	if hi.Settings.ShowComments {
		return hi.Continue()
	}

	switch hi.CurrentStatement.(type) {
	case *ast.Comment:
		if !hi.Settings.ShowComments {
			return hi.Stop()
		}
	}

	return hi.Continue()
}

// FilterDisabledStatements will filter out Assignment Statements that are
// disabled
func FilterDisabledStatements(hi *HandlerInput) HandlerSignal {
	// Short circuit the filter if we allow disabled statements
	if hi.Settings.IncludeDisabled {
		return hi.Continue()
	}

	switch statement := hi.CurrentStatement.(type) {
	case *ast.Assignment:
		if !statement.Active && !hi.Settings.IncludeDisabled {
			return hi.Stop()
		}
	}

	return hi.Continue()
}

// FilterByGroupName will filter out Statements that do not
// belong to the required Group name
func FilterByGroupName(hi *HandlerInput) HandlerSignal {
	// Short circuit the filter if there is no Group name to filter on
	if len(hi.Settings.FilterGroup) == 0 {
		return hi.Continue()
	}

	switch statement := hi.CurrentStatement.(type) {
	case *ast.Assignment:
		if !statement.BelongsToGroup(hi.Settings.FilterGroup) {
			return hi.Stop()
		}

	case *ast.Group:
		if !statement.BelongsToGroup(hi.Settings.FilterGroup) {
			return hi.Stop()
		}

	case *ast.Comment:
		if !statement.BelongsToGroup(hi.Settings.FilterGroup) {
			return hi.Stop()
		}
	}

	return hi.Continue()
}
