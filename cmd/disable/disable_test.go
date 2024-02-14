package disable_test

import (
	"testing"

	"github.com/jippi/dottie/pkg/test_helpers"
)

func TestDisable(t *testing.T) {
	t.Parallel()

	test_helpers.RunFilebasedCommandTests(t)
}
