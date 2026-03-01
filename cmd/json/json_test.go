package json_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/jippi/dottie/cmd"
	"github.com/jippi/dottie/pkg/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJsonCommandOutputsDocumentAsJSON(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, ".env")

	require.NoError(t, os.WriteFile(envFile, []byte("KEY_A=hello\nKEY_B=world\n"), 0o600))

	var stdout bytes.Buffer

	var stderr bytes.Buffer

	ctx := test_helpers.CreateTestContext(t, &stdout, &stderr)
	executed, err := cmd.RunCommand(ctx, []string{"json", "--file", envFile}, &stdout, &stderr)

	require.NoError(t, err)
	require.NotNil(t, executed)
	assert.Equal(t, "json", executed.Name())
	assert.Contains(t, stdout.String(), "\"KEY_A\"")
	assert.Contains(t, stdout.String(), "\"hello\"")
	assert.Contains(t, stdout.String(), "\"KEY_B\"")
	assert.Contains(t, stdout.String(), "\"world\"")
	assert.Empty(t, stderr.String())
}

func TestJsonCommandRejectsExtraArgs(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer

	var stderr bytes.Buffer

	ctx := test_helpers.CreateTestContext(t, &stdout, &stderr)
	executed, err := cmd.RunCommand(ctx, []string{"json", "extra-arg"}, &stdout, &stderr)

	require.Error(t, err)
	require.NotNil(t, executed)
	assert.Equal(t, "json", executed.Name())
	assert.Contains(t, err.Error(), "accepts 0 arg(s), received 1")
	assert.Contains(t, stderr.String(), "Run 'dottie json --help' for usage.")
	assert.Empty(t, stdout.String())
}

func TestJsonCommandErrorsWhenFileIsMissing(t *testing.T) {
	t.Parallel()

	missingFile := filepath.Join(t.TempDir(), "does-not-exist.env")

	var stdout bytes.Buffer

	var stderr bytes.Buffer

	ctx := test_helpers.CreateTestContext(t, &stdout, &stderr)
	executed, err := cmd.RunCommand(ctx, []string{"json", "--file", missingFile}, &stdout, &stderr)

	require.Error(t, err)
	require.NotNil(t, executed)
	assert.Equal(t, "json", executed.Name())
	assert.Contains(t, err.Error(), "no such file or directory")
	assert.Contains(t, stderr.String(), "Run 'dottie json --help' for usage.")
	assert.Empty(t, stdout.String())
}
