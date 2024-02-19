package ast

type selectorResult uint

const (
	Exclude selectorResult = iota
	Keep
)

type Selector func(input Statement) selectorResult
