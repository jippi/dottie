package print_cmd_test

import (
	"bytes"
	"testing"

	"github.com/jippi/dottie/cmd"
	"github.com/jippi/dottie/pkg/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrintCommand(t *testing.T) {
	t.Parallel()

	test_helpers.RunFileBasedCommandTests(t, test_helpers.ReadOnly, "print")
}

func TestPrintCommandNoColorEnv(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	ctx := test_helpers.CreateTestContext(t, &stdout, &stderr)

	_, err := cmd.RunCommand(
		ctx,
		[]string{"print", "--pretty", "--color", "--file", "tests/simple-color-explicit.env"},
		&stdout,
		&stderr,
	)

	require.NoError(t, err)
	assert.NotContains(t, stdout.String(), "\x1b[")
	assert.Contains(t, stdout.String(), "KEY_A=\"I'm key A\"")
}
