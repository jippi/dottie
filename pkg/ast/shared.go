package ast

// Statement represents syntax tree node of .env file statement (like: assignment or comment).
type Statement interface {
	statementNode()
	Is(Statement) bool
	BelongsToGroup(RenderSettings) bool
	ShouldRender(RenderSettings) bool
}
