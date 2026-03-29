package update_test

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

func TestSetCommand(t *testing.T) {
	t.Parallel()

	test_helpers.RunFileBasedCommandTests(t, 0, "update")
}

func TestCreateIfMissing(t *testing.T) {
	t.Parallel()

	targetFile := filepath.Join(t.TempDir(), "missing.env")

	var stdout, stderr bytes.Buffer

	ctx := test_helpers.CreateTestContext(t, &stdout, &stderr)

	_, err := cmd.RunCommand(ctx, []string{
		"update",
		"--source", "tests/source-flag.source",
		"--no-backup",
		"--no-validate",
		"--file", targetFile,
	}, &stdout, &stderr)
	require.NoError(t, err)

	content, err := os.ReadFile(targetFile)
	require.NoError(t, err)

	assert.Contains(t, string(content), `KEY="source"`)
}
