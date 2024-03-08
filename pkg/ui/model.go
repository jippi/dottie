package ui

import (
	"context"

	"github.com/jippi/dottie/pkg/ast"
	zone "github.com/lrstanley/bubblezone"
)

func NewModel(ctx context.Context, document *ast.Document) model {
	// Initialize a global zone manager, so we don't have to pass around the manager
	// throughout components.
	zone.NewGlobal()

	return model{
		document: document,
		form: form{
			ctx:      ctx,
			document: document,
		},
		groups: group{
			id:    zone.NewPrefix(),
			title: "Groups",
			items: Map(document.Groups, func(g *ast.Group) groupItem {
				return groupItem{name: g.String()}
			}),
		},
	}
}
