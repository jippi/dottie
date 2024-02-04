package ast

import (
	"reflect"
)

type Newline struct {
	Blank    bool     `json:"blank"`
	Group    *Group   `json:"-"`
	Position Position `json:"position"`
}

func (n *Newline) Is(other Statement) bool {
	return reflect.TypeOf(n) == reflect.TypeOf(other)
}

func (n *Newline) statementNode() {
}
