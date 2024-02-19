package render

import (
	"context"

	"github.com/jippi/dottie/pkg/ast"
)

func NewAstSelectorHandler(selectors ...ast.Selector) Handler {
	return func(ctx context.Context, input *HandlerInput) HandlerSignal {
		switch stmt := input.CurrentStatement.(type) {
		case ast.Statement:
			for _, selector := range selectors {
				if selector(stmt) == ast.Exclude {
					return input.Stop()
				}
			}
		}

		return input.Continue()
	}
}
