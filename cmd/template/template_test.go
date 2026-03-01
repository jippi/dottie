package template_test

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

type executionResult struct {
	stdout   string
	stderr   string
	executed string
	err      error
}

func runTemplateCommand(t *testing.T, envContent, templateContent string, args ...string) executionResult {
	t.Helper()

	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, ".env")
	templateFile := filepath.Join(tempDir, "tmpl.txt")

	require.NoError(t, os.WriteFile(envFile, []byte(envContent), 0o600))
	require.NoError(t, os.WriteFile(templateFile, []byte(templateContent), 0o600))

	var stdout bytes.Buffer

	var stderr bytes.Buffer

	commandArgs := append([]string{"template", templateFile}, args...)
	commandArgs = append(commandArgs, "--file", envFile)

	ctx := test_helpers.CreateTestContext(t, &stdout, &stderr)
	executed, err := cmd.RunCommand(ctx, commandArgs, &stdout, &stderr)

	result := executionResult{stdout: stdout.String(), stderr: stderr.String(), err: err}
	if executed != nil {
		result.executed = executed.Name()
	}

	return result
}

func runTemplateCommandWithMissingTemplateFile(t *testing.T) executionResult {
	t.Helper()

	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, ".env")
	missingTemplateFile := filepath.Join(tempDir, "missing-template.txt")

	require.NoError(t, os.WriteFile(envFile, []byte("A=hello\n"), 0o600))

	var stdout bytes.Buffer

	var stderr bytes.Buffer

	ctx := test_helpers.CreateTestContext(t, &stdout, &stderr)
	executed, err := cmd.RunCommand(ctx, []string{"template", missingTemplateFile, "--file", envFile}, &stdout, &stderr)

	result := executionResult{stdout: stdout.String(), stderr: stderr.String(), err: err}
	if executed != nil {
		result.executed = executed.Name()
	}

	return result
}

func TestTemplateCommandInterpolatesByDefault(t *testing.T) {
	t.Parallel()

	result := runTemplateCommand(t, "A=hello\nB=${A}\n", "{{ ( .Get \"B\" ).Interpolated }}")

	require.NoError(t, result.err)
	assert.Equal(t, "template", result.executed)
	assert.Equal(t, "hello", result.stdout)
	assert.Empty(t, result.stderr)
}

func TestTemplateCommandNoInterpolationFlag(t *testing.T) {
	t.Parallel()

	result := runTemplateCommand(t, "A=hello\nB=${A}\n", "{{ ( .Get \"B\" ).Interpolated }}", "--no-interpolation")

	require.NoError(t, result.err)
	assert.Equal(t, "template", result.executed)
	assert.Equal(t, "hello", result.stdout)
	assert.Empty(t, result.stderr)
}

func TestTemplateCommandWithDisabledFlagInterpolatesDisabledAssignments(t *testing.T) {
	t.Parallel()

	withoutFlag := runTemplateCommand(t, "B=hello\n#A=${B}\n", "{{ ( .Get \"A\" ).Interpolated }}")
	withFlag := runTemplateCommand(t, "B=hello\n#A=${B}\n", "{{ ( .Get \"A\" ).Interpolated }}", "--with-disabled")

	require.NoError(t, withoutFlag.err)
	assert.Empty(t, withoutFlag.stdout)

	require.NoError(t, withFlag.err)
	assert.Equal(t, "template", withFlag.executed)
	assert.Equal(t, "hello", withFlag.stdout)
	assert.Empty(t, withFlag.stderr)
}

func TestTemplateCommandErrorsWhenTemplateFileIsMissing(t *testing.T) {
	t.Parallel()

	result := runTemplateCommandWithMissingTemplateFile(t)

	require.Error(t, result.err)
	assert.Equal(t, "template", result.executed)
	assert.Contains(t, result.err.Error(), "no such file or directory")
	assert.Contains(t, result.stderr, "Run 'dottie template --help' for usage.")
	assert.Empty(t, result.stdout)
}

func TestTemplateCommandRejectsMissingAndExtraArgs(t *testing.T) {
	t.Parallel()

	t.Run("missing template file arg", func(t *testing.T) {
		t.Parallel()

		var stdout bytes.Buffer

		var stderr bytes.Buffer

		ctx := test_helpers.CreateTestContext(t, &stdout, &stderr)
		executed, err := cmd.RunCommand(ctx, []string{"template"}, &stdout, &stderr)

		require.Error(t, err)
		require.NotNil(t, executed)
		assert.Equal(t, "template", executed.Name())
		assert.Contains(t, err.Error(), "accepts 1 arg(s), received 0")
		assert.Contains(t, stderr.String(), "Run 'dottie template --help' for usage.")
		assert.Empty(t, stdout.String())
	})

	t.Run("extra args", func(t *testing.T) {
		t.Parallel()

		var stdout bytes.Buffer

		var stderr bytes.Buffer

		ctx := test_helpers.CreateTestContext(t, &stdout, &stderr)
		executed, err := cmd.RunCommand(ctx, []string{"template", "one", "two"}, &stdout, &stderr)

		require.Error(t, err)
		require.NotNil(t, executed)
		assert.Equal(t, "template", executed.Name())
		assert.Contains(t, err.Error(), "accepts 1 arg(s), received 2")
		assert.Contains(t, stderr.String(), "Run 'dottie template --help' for usage.")
		assert.Empty(t, stdout.String())
	})
}

func TestTemplateCommandPanicsOnInvalidTemplateSyntax(t *testing.T) {
	t.Parallel()

	require.Panics(t, func() {
		_ = runTemplateCommand(t, "A=hello\n", "{{ if }}")
	})
}

func TestTemplateCommandReturnsExecutionError(t *testing.T) {
	t.Parallel()

	result := runTemplateCommand(t, "A=hello\n", "{{ index . \"MISSING\" }}")

	require.Error(t, result.err)
	assert.Equal(t, "template", result.executed)
	assert.Contains(t, result.err.Error(), "can't index item")
	assert.Contains(t, result.stderr, "Run 'dottie template --help' for usage.")
	assert.Empty(t, result.stdout)
}
