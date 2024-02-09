package render

import (
	"strings"

	"github.com/jippi/dottie/pkg/ast"
)

// FilterComments will filter out Comment statements if they aren't to be included
func FilterComments(input *HandlerInput) HandlerSignal {
	// Short circuit the filter if we allow comments
	if input.Settings.showComments {
		return input.Continue()
	}

	switch input.CurrentStatement.(type) {
	case *ast.Comment:
		if !input.Settings.showComments {
			return input.Stop()
		}
	}

	return input.Continue()
}

// FilterDisabledStatements will filter out Assignment Statements that are
// disabled
func FilterDisabledStatements(input *HandlerInput) HandlerSignal {
	// Short circuit the filter if we allow disabled statements
	if input.Settings.includeDisabled {
		return input.Continue()
	}

	switch statement := input.CurrentStatement.(type) {
	case *ast.Assignment:
		if !statement.Active && !input.Settings.includeDisabled {
			return input.Stop()
		}
	}

	return input.Continue()
}

// FilterActiveStatements will filter out Assignment Statements that are
// *active*
func FilterActiveStatements(input *HandlerInput) HandlerSignal {
	switch statement := input.CurrentStatement.(type) {
	case *ast.Assignment:
		if statement.Active {
			return input.Stop()
		}
	}

	return input.Continue()
}

// FilterGroupName will filter out Statements that do not
// belong to the required Group name
func FilterGroupName(input *HandlerInput) HandlerSignal {
	// Short circuit the filter if there is no Group name to filter on
	if len(input.Settings.filterGroup) == 0 {
		return input.Continue()
	}

	switch statement := input.CurrentStatement.(type) {
	case *ast.Assignment:
		if !statement.BelongsToGroup(input.Settings.filterGroup) {
			return input.Stop()
		}

	case *ast.Group:
		if !statement.BelongsToGroup(input.Settings.filterGroup) {
			return input.Stop()
		}

	case *ast.Comment:
		if !statement.BelongsToGroup(input.Settings.filterGroup) {
			return input.Stop()
		}
	}

	return input.Continue()
}

// FilterKeyPrefix will filter out Assignment Statements that do not have the
// configured (optional) key prefix.
func FilterKeyPrefix(input *HandlerInput) HandlerSignal {
	// Short circuit the filter if there is no KeyPrefix to filter on
	if len(input.Settings.filterKeyPrefix) == 0 {
		return input.Continue()
	}

	switch statement := input.CurrentStatement.(type) {
	case *ast.Assignment:
		if !strings.HasPrefix(statement.Name, input.Settings.filterKeyPrefix) {
			return input.Stop()
		}
	}

	return input.Continue()
}
