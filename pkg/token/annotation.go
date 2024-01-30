package token

import (
	"fmt"
	"strings"
)

type Annotation struct {
	Key   string
	Value string
}

func (a Annotation) String() string {
	return fmt.Sprintf("ANNOTATION(KEY=%s;VALUE=%s)", a.Key, a.Value)
}

func (a Annotation) IsDottie() bool {
	return strings.HasPrefix(a.Key, "dottie/")
}
