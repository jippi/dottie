package ast

// RetainGroup will exclude Assignment, Group, and Comment statements
// that do not belong to the provided group name
func RetainGroup(name string) Selector {
	return func(input Statement) selectorResult {
		switch statement := input.(type) {
		case *Assignment:
			if !statement.BelongsToGroup(name) {
				return Exclude
			}

		case *Group:
			if !statement.BelongsToGroup(name) {
				return Exclude
			}

		case *Comment:
			if !statement.BelongsToGroup(name) {
				return Exclude
			}
		}

		return Keep
	}
}
