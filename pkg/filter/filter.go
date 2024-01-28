package filter

import (
	"strings"

	"dotfedi/pkg/ast"
)

type Filter struct {
	KeyPrefix string
	Group     string
	Commented bool
}

func (f *Filter) Match(assignment *ast.Assignment) bool {
	if len(f.KeyPrefix) > 0 && !strings.HasPrefix(assignment.Key, f.KeyPrefix) {
		return false
	}

	if len(f.Group) > 0 && assignment.Group != nil && assignment.Group.Comment != f.Group {
		return false
	}

	if assignment.Commented && !f.Commented {
		return false
	}

	return true
}
