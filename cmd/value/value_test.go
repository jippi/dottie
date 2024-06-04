package value_test

import (
	"testing"

	"github.com/jippi/dottie/pkg/test_helpers"
)

func TestPrintCommand(t *testing.T) {
	t.Parallel()

	test_helpers.RunFileBasedCommandTests(t, test_helpers.ReadOnly, "value")
}
