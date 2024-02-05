package render

import (
	"github.com/jippi/dottie/pkg/ast"
)

type Handler func(in *HandlerInput) HandlerSignal

type HandlerInput struct {
	Presenter         *Renderer
	PreviousStatement ast.Statement
	Settings          Settings
	CurrentStatement  any
	ReturnValue       *LineBuffer
}

func (hi *HandlerInput) Stop() HandlerSignal {
	return Stop
}

func (hi *HandlerInput) Return(value *LineBuffer) HandlerSignal {
	hi.ReturnValue = value

	return Return
}

func (hi *HandlerInput) Continue() HandlerSignal {
	return Continue
}
