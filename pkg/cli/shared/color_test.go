package shared_test

import (
	"testing"

	"github.com/jippi/dottie/pkg/cli/shared"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

func TestColorEnabled(t *testing.T) {
	t.Parallel()

	t.Run("default true", func(t *testing.T) {
		t.Parallel()

		flags := testColorFlags(t)
		require.True(t, shared.ColorEnabled(flags, "color"))
	})

	t.Run("explicit no-color", func(t *testing.T) {
		t.Parallel()

		flags := testColorFlags(t)
		require.NoError(t, flags.Set("no-color", "true"))
		require.False(t, shared.ColorEnabled(flags, "color"))
	})

	t.Run("explicit color false", func(t *testing.T) {
		t.Parallel()

		flags := testColorFlags(t)
		require.NoError(t, flags.Set("color", "false"))
		require.False(t, shared.ColorEnabled(flags, "color"))
	})
}

func TestColorEnabledNoColorEnvVar(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	flags := testColorFlags(t)
	require.NoError(t, flags.Set("color", "true"))
	require.False(t, shared.ColorEnabled(flags, "color"))
}

func testColorFlags(t *testing.T) *pflag.FlagSet {
	t.Helper()

	cmd := &cobra.Command{Use: "test"}
	shared.BoolWithInverse(cmd, "color", true, "", "")

	return cmd.Flags()
}
