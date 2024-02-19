package upsert

import "fmt"

type SkippedStatementError struct {
	Key    string
	Reason string
}

func (e SkippedStatementError) Error() string {
	return fmt.Sprintf("Key [ %s ] was skipped: %s", e.Key, e.Reason)
}
