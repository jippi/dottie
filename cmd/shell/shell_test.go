package shell_test

import (
	"testing"

	"github.com/jippi/dottie/pkg/test_helpers"
)

func TestShellCommand(t *testing.T) {
	test_helpers.RunFileBasedCommandTests(t, 0, "shell")
}
