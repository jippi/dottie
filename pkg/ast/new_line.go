package ast

import (
	"reflect"
)

type Newline struct {
	Blank    bool     `json:"blank"`
	Group    *Group   `json:"-"`
	Repeated int      `json:"repeated"`
	Position Position `json:"position"`
}

func (n *Newline) Is(other Statement) bool {
	if n == nil || other == nil {
		return false
	}

	return n.Type() == other.Type()
}

func (n *Newline) Type() string {
	if n == nil {
		return "<nil>Newline"
	}

	return reflect.TypeOf(n).String()
}

func (n *Newline) statementNode() {
}
