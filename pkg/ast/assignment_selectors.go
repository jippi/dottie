package ast

import (
	"strings"
)

// ExcludeHiddenViaAnnotation will exclude *HIDDEN* Assignments
// via the [@dottie/hidden] annotation
func ExcludeHiddenViaAnnotation(input Statement) SelectorResult {
	switch statement := input.(type) {
	case *Assignment:
		if statement.IsHidden() {
			return Exclude
		}
	}

	return Keep
}

// ExcludeDisabledAssignments will exclude *DISABLED* Assignments
func ExcludeDisabledAssignments(input Statement) SelectorResult {
	switch statement := input.(type) {
	case *Assignment:
		if !statement.Active {
			return Exclude
		}
	}

	return Keep
}

// ExcludeActiveAssignments will exclude *ACTIVE* Assignments
func ExcludeActiveAssignments(input Statement) SelectorResult {
	switch statement := input.(type) {
	case *Assignment:
		if statement.Active {
			return Exclude
		}
	}

	return Keep
}

// RetainKeyPrefix will *RETAIN* Assignments with the provided prefix
func RetainKeyPrefix(prefix string) Selector {
	return func(input Statement) SelectorResult {
		switch statement := input.(type) {
		case *Assignment:
			if !strings.HasPrefix(statement.Name, prefix) {
				return Exclude
			}
		}

		return Keep
	}
}

// ExcludeKeyPrefix will *EXCLUDE* Assignments with the provided prefix
func ExcludeKeyPrefix(prefix string) Selector {
	return func(input Statement) SelectorResult {
		switch statement := input.(type) {
		case *Assignment:
			if strings.HasPrefix(statement.Name, prefix) {
				return Exclude
			}
		}

		return Keep
	}
}
