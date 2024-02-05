package render

import (
	"github.com/jippi/dottie/pkg/ast"
)

type Handler func(hi *HandlerInput) HandlerSignal

type HandlerInput struct {
	CurrentStatement  any
	PreviousStatement ast.Statement
	Renderer          *Renderer
	ReturnValue       *Lines
	Settings          Settings
}

func (hi *HandlerInput) Stop() HandlerSignal {
	return Stop
}

func (hi *HandlerInput) Return(value *Lines) HandlerSignal {
	hi.ReturnValue = value

	return Return
}

func (hi *HandlerInput) Continue() HandlerSignal {
	return Continue
}
