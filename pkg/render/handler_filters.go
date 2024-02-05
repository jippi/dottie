package render

import (
	"strings"

	"github.com/jippi/dottie/pkg/ast"
)

// FilterComments will filter out Comment statements if they aren't to be included
func FilterComments(hi *HandlerInput) HandlerSignal {
	// Short circuit the filter if we allow comments
	if hi.Settings.showComments {
		return hi.Continue()
	}

	switch hi.CurrentStatement.(type) {
	case *ast.Comment:
		if !hi.Settings.showComments {
			return hi.Stop()
		}
	}

	return hi.Continue()
}

// FilterDisabledStatements will filter out Assignment Statements that are
// disabled
func FilterDisabledStatements(hi *HandlerInput) HandlerSignal {
	// Short circuit the filter if we allow disabled statements
	if hi.Settings.includeDisabled {
		return hi.Continue()
	}

	switch statement := hi.CurrentStatement.(type) {
	case *ast.Assignment:
		if !statement.Active && !hi.Settings.includeDisabled {
			return hi.Stop()
		}
	}

	return hi.Continue()
}

// FilterGroupName will filter out Statements that do not
// belong to the required Group name
func FilterGroupName(hi *HandlerInput) HandlerSignal {
	// Short circuit the filter if there is no Group name to filter on
	if len(hi.Settings.filterGroup) == 0 {
		return hi.Continue()
	}

	switch statement := hi.CurrentStatement.(type) {
	case *ast.Assignment:
		if !statement.BelongsToGroup(hi.Settings.filterGroup) {
			return hi.Stop()
		}

	case *ast.Group:
		if !statement.BelongsToGroup(hi.Settings.filterGroup) {
			return hi.Stop()
		}

	case *ast.Comment:
		if !statement.BelongsToGroup(hi.Settings.filterGroup) {
			return hi.Stop()
		}
	}

	return hi.Continue()
}

// FilterKeyPrefix will filter out Assignment Statements that do not have the
// configured (optional) key prefix.
func FilterKeyPrefix(hi *HandlerInput) HandlerSignal {
	// Short circuit the filter if there is no KeyPrefix to filter on
	if len(hi.Settings.filterKeyPrefix) == 0 {
		return hi.Continue()
	}

	switch statement := hi.CurrentStatement.(type) {
	case *ast.Assignment:
		if !strings.HasPrefix(statement.Name, hi.Settings.filterKeyPrefix) {
			return hi.Stop()
		}
	}

	return hi.Continue()
}
