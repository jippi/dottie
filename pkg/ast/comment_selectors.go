package ast

// ExcludeComments will *EXCLUDE* all comments
func ExcludeComments(input Statement) selectorResult {
	switch input.(type) {
	case *Comment:
		return Exclude
	}

	return Keep
}
