package fmt_test

import (
	"testing"

	"github.com/jippi/dottie/pkg/test_helpers"
)

func TestEnableCommand(t *testing.T) {
	t.Parallel()

	test_helpers.RunFileBasedCommandTests(t, 0, "fmt")
}
