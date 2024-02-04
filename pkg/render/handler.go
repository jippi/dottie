package render

import (
	"github.com/jippi/dottie/pkg/ast"
)

type Handler func(in *HandlerInput) HandlerSignal

type HandlerInput struct {
	Presenter *Renderer
	Previous  ast.Statement
	Settings  Settings
	Statement any
	Value     string
}

func (si *HandlerInput) Stop() HandlerSignal {
	return Stop
}

func (si *HandlerInput) Return(val string) HandlerSignal {
	si.Value = val

	return Return
}

func (si *HandlerInput) Continue() HandlerSignal {
	return Continue
}
