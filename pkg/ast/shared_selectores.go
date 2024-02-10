package ast

type SelectorResult uint

const (
	Exclude SelectorResult = iota
	Keep
)

type Selector func(input Statement) SelectorResult
