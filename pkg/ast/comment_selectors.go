package ast

// ExcludeComments will *EXCLUDE* all comments
func ExcludeComments(input Statement) SelectorResult {
	switch input.(type) {
	case *Comment:
		return Exclude
	}

	return Keep
}
