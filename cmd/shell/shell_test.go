package shell_test

import (
	"bytes"
	"testing"

	"github.com/jippi/dottie/cmd/shell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShellCommandMetadata(t *testing.T) {
	t.Parallel()

	command := shell.New()

	assert.Equal(t, "shell", command.Use)
	assert.Equal(t, "manipulate", command.GroupID)
	require.NotNil(t, command.Args)
	require.NotNil(t, command.RunE)
}

func TestShellCommandArgsValidation(t *testing.T) {
	t.Parallel()

	command := shell.New()

	require.NoError(t, command.Args(command, []string{}))

	err := command.Args(command, []string{"unexpected"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "accepts 0 arg(s), received 1")
}

func TestShellCommandExecuteRejectsExtraArgs(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer

	var stderr bytes.Buffer

	command := shell.New()
	command.SetOut(&stdout)
	command.SetErr(&stderr)
	command.SetArgs([]string{"unexpected"})

	err := command.Execute()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "accepts 0 arg(s), received 1")
	assert.Equal(t, "shell", command.Name())
	assert.Contains(t, stderr.String(), "accepts 0 arg(s), received 1")
}
