package update_test

import (
	"testing"

	"github.com/jippi/dottie/pkg/test_helpers"
)

func TestSetCommand(t *testing.T) {
	t.Parallel()

	test_helpers.RunFileBasedCommandTests(t, 0, "update")
}
