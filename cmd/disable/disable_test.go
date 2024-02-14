package disable_test

import (
	"testing"

	"github.com/jippi/dottie/pkg/test_helpers"
)

func TestCommand(t *testing.T) {
	t.Parallel()

	test_helpers.RunFilebasedCommandTests(t, "disable")
}
