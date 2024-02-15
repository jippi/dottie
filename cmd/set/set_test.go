package set_test

import (
	"testing"

	"github.com/jippi/dottie/pkg/test_helpers"
)

func TestCommand(t *testing.T) {
	t.Parallel()

	test_helpers.RunFileBasedCommandTests(t, 0, "set")
}
