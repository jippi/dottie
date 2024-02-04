package render

import "github.com/jippi/dottie/pkg/ast"

type Presenter interface {
	SetOutput(output Outputter)
	Statement(stmt any, previous ast.Statement, settings Settings) string
}
