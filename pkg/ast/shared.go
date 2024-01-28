package ast

// Statement represents syntax tree node of .env file statement (like: assignment or comment).
type Statement interface {
	statementNode()
	Is(Statement) bool
	BelongsToGroup(RenderSettings) bool
	ShouldRender(RenderSettings) bool
}

// Type is the set of lexical tokens.
type Type uint

// The list of tokens.
const (
	Illegal Type = iota
)
