package render

import (
	"github.com/jippi/dottie/pkg/ast"
)

type Signal uint

const (
	Continue Signal = iota
	Stop
	Return
)

var signals = []string{
	Continue: "CONTINUE",
	Stop:     "STOP",
	Return:   "RETURN",
}

// String returns the string corresponding to the token.
func (ss Signal) String() string {
	s := ""

	if int(ss) < len(signals) {
		s = signals[ss]
	}

	return s
}

type Handler func(in *HandlerInput) Signal

type HandlerInput struct {
	Presenter *Renderer
	Previous  ast.Statement
	Settings  Settings
	Statement any
	Value     string
}

func (si *HandlerInput) Stop() Signal {
	return Stop
}

func (si *HandlerInput) Return(val string) Signal {
	si.Value = val

	return Return
}

func (si *HandlerInput) Continue() Signal {
	return Continue
}
