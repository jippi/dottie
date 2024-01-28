package filter

import (
	"strings"

	"dotfedi/pkg/ast"
)

type Filter struct {
	KeyPrefix string
	Group     string
}

func (f *Filter) Match(assignment *ast.Assignment) bool {
	if len(f.KeyPrefix) > 0 && !strings.HasPrefix(assignment.Key, f.KeyPrefix) {
		return false
	}

	if len(f.Group) > 0 && assignment.Group != nil && assignment.Group.Comment != f.Group {
		return false
	}

	return true
}
