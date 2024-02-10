package render

import (
	"github.com/jippi/dottie/pkg/ast"
)

var (
	ExcludeComments            = newSelectorHandler(ast.ExcludeComments)
	ExcludeDisabledAssignments = newSelectorHandler(ast.ExcludeDisabledAssignments)
	ExcludeActiveAssignments   = newSelectorHandler(ast.ExcludeActiveAssignments)
	ExcludeHiddenViaAnnotation = newSelectorHandler(ast.ExcludeHiddenViaAnnotation)
)

func RetainGroup(value string) Handler       { return newSelectorHandler(ast.RetainGroup(value)) }
func ExcludeKeyPrefix(value string) Handler  { return newSelectorHandler(ast.ExcludeKeyPrefix(value)) }
func RetainKeyPrefix(value string) Handler   { return newSelectorHandler(ast.RetainKeyPrefix(value)) }
func RetainExactKey(value ...string) Handler { return newSelectorHandler(ast.RetainExactKey(value...)) }

func newSelectorHandler(selector ast.Selector) Handler {
	return func(input *HandlerInput) HandlerSignal {
		switch stmt := input.CurrentStatement.(type) {
		case ast.Statement:
			if selector(stmt) == ast.Exclude {
				return input.Stop()
			}
		}

		return input.Continue()
	}
}
